package user_test

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ivankornilov/auth/internal/model"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

func (s *APISuite) TestGet() {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	tests := []struct {
		name    string
		id      int64
		user    *model.User
		service error
		want    *desc.GetResponse
		code    codes.Code
		message string
	}{
		{
			name: "success",
			id:   7,
			user: &model.User{
				ID: 7,
				Info: model.UserInfo{
					Name:  "Ivan",
					Email: "ivan@example.com",
					Role:  "ADMIN",
				},
				CreatedAt: createdAt,
				UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
			},
			want: &desc.GetResponse{
				Id:        7,
				Name:      "Ivan",
				Email:     "ivan@example.com",
				Role:      desc.Role_ADMIN,
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		},
		{
			name: "without updated_at",
			id:   7,
			user: &model.User{
				ID: 7,
				Info: model.UserInfo{
					Name:  "Ivan",
					Email: "ivan@example.com",
					Role:  "USER",
				},
				CreatedAt: createdAt,
			},
			want: &desc.GetResponse{
				Id:        7,
				Name:      "Ivan",
				Email:     "ivan@example.com",
				Role:      desc.Role_USER,
				CreatedAt: timestamppb.New(createdAt),
			},
		},
		{
			name: "unknown role",
			id:   7,
			user: &model.User{
				ID: 7,
				Info: model.UserInfo{
					Name:  "Ivan",
					Email: "ivan@example.com",
					Role:  "BOSS",
				},
				CreatedAt: createdAt,
			},
			want: &desc.GetResponse{
				Id:        7,
				Name:      "Ivan",
				Email:     "ivan@example.com",
				Role:      desc.Role_USER,
				CreatedAt: timestamppb.New(createdAt),
			},
		},
		{
			name:    "not found",
			id:      7,
			service: model.ErrUserNotFound,
			code:    codes.NotFound,
			message: fmt.Sprintf("user with id %d not found", 7),
		},
		{
			name:    "internal",
			id:      7,
			service: errors.New("db is down"),
			code:    codes.Internal,
			message: "db is down",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.setup()
			s.userService.GetMock.Expect(s.ctx, tt.id).Return(tt.user, tt.service)

			resp, err := s.api.Get(s.ctx, &desc.GetRequest{Id: tt.id})
			if tt.service != nil {
				s.Require().Nil(resp)
				s.requireStatus(err, tt.code, tt.message)
				return
			}

			s.Require().NoError(err)
			s.True(proto.Equal(tt.want, resp))
		})
	}
}
