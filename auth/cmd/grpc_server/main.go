package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ivankornilov/auth/internal/config"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetName() == "" || req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "name, email and password are required")
	}
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Error(codes.InvalidArgument, "password and password_confirm do not match")
	}

	query, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "role").
		Values(req.GetName(), req.GetEmail(), req.GetPassword(), req.GetRole().String()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	var id int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	log.Printf("created user id=%d", id)

	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	query, args, err := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	var (
		id              int64
		name, email     string
		roleStr         string
		createdAt       time.Time
		updatedAt       sql.NullTime
	)

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &roleStr, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	role := desc.Role_USER
	if value, ok := desc.Role_value[roleStr]; ok {
		role = desc.Role(value)
	}

	resp := &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: timestamppb.New(createdAt),
	}
	if updatedAt.Valid {
		resp.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return resp, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if req.GetName() == nil && req.GetEmail() == nil {
		return nil, status.Error(codes.InvalidArgument, "name or email is required")
	}

	builder := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.GetId()})

	if req.GetName() != nil {
		builder = builder.Set("name", req.GetName().GetValue())
	}
	if req.GetEmail() != nil {
		builder = builder.Set("email", req.GetEmail().GetValue())
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	tag, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	query, args, err := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	tag, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
	}

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
