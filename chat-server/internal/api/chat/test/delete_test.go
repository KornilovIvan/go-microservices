package chat_test

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ivankornilov/chat-server/internal/model"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

func (s *APISuite) TestDelete() {
	const id int64 = 4

	tests := []struct {
		name    string
		service error
		code    codes.Code
		message string
	}{
		{
			name: "success",
		},
		{
			name:    "not found",
			service: model.ErrChatNotFound,
			code:    codes.NotFound,
			message: fmt.Sprintf("chat with id %d not found", id),
		},
		{
			name:    "internal",
			service: errors.New("db is down"),
			code:    codes.Internal,
			message: "db is down",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.setup()
			s.chatService.DeleteMock.Expect(s.ctx, id).Return(tt.service)

			resp, err := s.api.Delete(s.ctx, &desc.DeleteRequest{Id: id})
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
