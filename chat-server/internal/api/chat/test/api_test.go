package chat_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivankornilov/chat-server/internal/api/chat"
	"github.com/ivankornilov/chat-server/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx         context.Context
	chatService *mocks.ChatServiceMock
	api         *chat.Implementation
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}

func (s *APISuite) setup() {
	s.ctx = context.Background()
	s.chatService = mocks.NewChatServiceMock(s.T())
	s.api = chat.NewImplementation(s.chatService)
}

func (s *APISuite) requireStatus(err error, code codes.Code, message string) {
	s.T().Helper()
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(code, st.Code())
	s.Equal(message, st.Message())
}
