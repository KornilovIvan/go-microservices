package chat_test

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

func (s *APISuite) TestCreate() {
	usernames := []string{"alice", "bob"}
	req := &desc.CreateRequest{Usernames: usernames}

	tests := []struct {
		name    string
		service error
		id      int64
		code    codes.Code
		message string
	}{
		{
			name: "success",
			id:   3,
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
			s.chatService.CreateMock.Expect(s.ctx, usernames).Return(tt.id, tt.service)

			resp, err := s.api.Create(s.ctx, req)
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
