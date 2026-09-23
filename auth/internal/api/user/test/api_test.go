package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivankornilov/auth/internal/api/user"
	"github.com/ivankornilov/auth/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx         context.Context
	userService *mocks.UserServiceMock
	api         *user.Implementation
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}

func (s *APISuite) setup() {
	s.ctx = context.Background()
	s.userService = mocks.NewUserServiceMock(s.T())
	s.api = user.NewImplementation(s.userService, nil)
}

func (s *APISuite) requireStatus(err error, code codes.Code, message string) {
	s.T().Helper()
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(code, st.Code())
	s.Equal(message, st.Message())
}
