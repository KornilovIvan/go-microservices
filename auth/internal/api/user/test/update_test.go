package user_test

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/ivankornilov/auth/internal/model"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

func (s *APISuite) TestUpdate() {
	name := "Ivan"
	email := "ivan@example.com"

	tests := []struct {
		name    string
		req     *desc.UpdateRequest
		expect  *model.UpdateUser
		service error
		code    codes.Code
		message string
	}{
		{
			name: "success",
			req: &desc.UpdateRequest{
				Id:    7,
				Name:  wrapperspb.String(name),
				Email: wrapperspb.String(email),
			},
			expect: &model.UpdateUser{ID: 7, Name: &name, Email: &email},
		},
		{
			name: "name only",
			req: &desc.UpdateRequest{
				Id:   7,
				Name: wrapperspb.String(name),
			},
			expect: &model.UpdateUser{ID: 7, Name: &name},
		},
		{
			name:    "name or email required",
			req:     &desc.UpdateRequest{Id: 7},
			expect:  &model.UpdateUser{ID: 7},
			service: model.ErrNameOrEmailRequired,
			code:    codes.InvalidArgument,
			message: model.ErrNameOrEmailRequired.Error(),
		},
		{
			name: "not found",
			req: &desc.UpdateRequest{
				Id:   7,
				Name: wrapperspb.String(name),
			},
			expect:  &model.UpdateUser{ID: 7, Name: &name},
			service: model.ErrUserNotFound,
			code:    codes.NotFound,
			message: fmt.Sprintf("user with id %d not found", 7),
		},
		{
			name: "already exists",
			req: &desc.UpdateRequest{
				Id:    7,
				Email: wrapperspb.String(email),
			},
			expect:  &model.UpdateUser{ID: 7, Email: &email},
			service: model.ErrUserAlreadyExists,
			code:    codes.AlreadyExists,
			message: model.ErrUserAlreadyExists.Error(),
		},
		{
			name: "internal",
			req: &desc.UpdateRequest{
				Id:   7,
				Name: wrapperspb.String(name),
			},
			expect:  &model.UpdateUser{ID: 7, Name: &name},
			service: errors.New("db is down"),
			code:    codes.Internal,
			message: "db is down",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.setup()
			s.userService.UpdateMock.Expect(s.ctx, tt.expect).Return(tt.service)

			resp, err := s.api.Update(s.ctx, tt.req)
			if tt.service != nil {
				s.Require().Nil(resp)
				s.requireStatus(err, tt.code, tt.message)
				return
			}

			s.Require().NoError(err)
			s.True(proto.Equal(&emptypb.Empty{}, resp))
		})
	}
}
