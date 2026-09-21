package user_test

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	"github.com/ivankornilov/auth/internal/model"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

func (s *APISuite) TestCreate() {
	info := &model.UserInfo{
		Name:            "Ivan",
		Email:           "ivan@example.com",
		Password:        "secret",
		PasswordConfirm: "secret",
		Role:            "USER",
	}
	req := &desc.CreateRequest{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            desc.Role_USER,
	}

	tests := []struct {
		name    string
		req     *desc.CreateRequest
		expect  *model.UserInfo
		service error
		id      int64
		code    codes.Code
		message string
	}{
		{
			name:   "success",
			req:    req,
			expect: info,
			id:     10,
		},
		{
			name: "required fields",
			req:  &desc.CreateRequest{Role: desc.Role_USER},
			expect: &model.UserInfo{
				Role: "USER",
			},
			service: model.ErrNameEmailPasswordRequired,
			code:    codes.InvalidArgument,
			message: model.ErrNameEmailPasswordRequired.Error(),
		},
		{
			name: "password mismatch",
			req: &desc.CreateRequest{
				Name:            "Ivan",
				Email:           "ivan@example.com",
				Password:        "secret",
				PasswordConfirm: "other",
				Role:            desc.Role_USER,
			},
			expect: &model.UserInfo{
				Name:            "Ivan",
				Email:           "ivan@example.com",
				Password:        "secret",
				PasswordConfirm: "other",
				Role:            "USER",
			},
			service: model.ErrPasswordMismatch,
			code:    codes.InvalidArgument,
			message: model.ErrPasswordMismatch.Error(),
		},
		{
			name:    "already exists",
			req:     req,
			expect:  info,
			service: model.ErrUserAlreadyExists,
			code:    codes.AlreadyExists,
			message: model.ErrUserAlreadyExists.Error(),
		},
		{
			name:    "internal",
			req:     req,
			expect:  info,
			service: errors.New("db is down"),
			code:    codes.Internal,
			message: "db is down",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.setup()
			s.userService.CreateMock.Expect(s.ctx, tt.expect).Return(tt.id, tt.service)

			resp, err := s.api.Create(s.ctx, tt.req)
			if tt.service != nil {
				s.Require().Nil(resp)
				s.requireStatus(err, tt.code, tt.message)
				return
			}

			s.Require().NoError(err)
			s.True(proto.Equal(&desc.CreateResponse{Id: tt.id}, resp))
		})
	}
}
