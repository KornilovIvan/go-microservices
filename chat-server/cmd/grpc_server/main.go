package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ivankornilov/chat-server/internal/config"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query, args, err := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("created_at").
		Values(time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	var chatID int64
	err = tx.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create chat: %v", err)
	}

	for _, username := range req.GetUsernames() {
		if username == "" {
			continue
		}

		userQuery, userArgs, err := sq.Insert("chat_users").
			PlaceholderFormat(sq.Dollar).
			Columns("chat_id", "username").
			Values(chatID, username).
			ToSql()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
		}

		_, err = tx.Exec(ctx, userQuery, userArgs...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to add chat user: %v", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	log.Printf("created chat id=%d usernames=%v", chatID, req.GetUsernames())

	return &desc.CreateResponse{Id: chatID}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	query, args, err := sq.Delete("chats").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	tag, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete chat: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "chat with id %d not found", req.GetId())
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if req.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat_id is required")
	}
	if req.GetFrom() == "" || req.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "from and text are required")
	}

	sentAt := time.Now()
	if req.GetTimestamp() != nil && req.GetTimestamp().IsValid() {
		sentAt = req.GetTimestamp().AsTime()
	}

	query, args, err := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "from_user", "text", "sent_at").
		Values(req.GetChatId(), req.GetFrom(), req.GetText(), sentAt).
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, status.Errorf(codes.NotFound, "chat with id %d not found", req.GetChatId())
		}
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
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
	desc.RegisterChatV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
