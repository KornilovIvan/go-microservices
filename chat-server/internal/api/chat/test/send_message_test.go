package chat_test

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ivankornilov/chat-server/internal/model"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

func (s *APISuite) TestSendMessage() {
	sentAt := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	req := &desc.SendMessageRequest{
		ChatId:    5,
		From:      "ivan",
		Text:      "hello",
		Timestamp: timestamppb.New(sentAt),
	}
	message := &model.Message{
		ChatID: req.ChatId,
		From:   req.From,
		Text:   req.Text,
		SentAt: req.GetTimestamp().AsTime(),
	}

	tests := []struct {
		name    string
		req     *desc.SendMessageRequest
		expect  *model.Message
		service error
		code    codes.Code
		message string
	}{
		{
			name:   "success",
			req:    req,
			expect: message,
		},
		{
			name: "chat id required",
			req: &desc.SendMessageRequest{
				From:      "ivan",
				Text:      "hello",
				Timestamp: timestamppb.New(sentAt),
			},
			expect: &model.Message{
				From:   "ivan",
				Text:   "hello",
				SentAt: req.GetTimestamp().AsTime(),
			},
			service: model.ErrChatIDRequired,
			code:    codes.InvalidArgument,
			message: model.ErrChatIDRequired.Error(),
		},
		{
			name: "from and text required",
			req: &desc.SendMessageRequest{
				ChatId:    5,
				Timestamp: timestamppb.New(sentAt),
			},
			expect: &model.Message{
				ChatID: 5,
				SentAt: req.GetTimestamp().AsTime(),
			},
			service: model.ErrFromAndTextRequired,
			code:    codes.InvalidArgument,
			message: model.ErrFromAndTextRequired.Error(),
		},
		{
			name:    "not found",
			req:     req,
			expect:  message,
			service: model.ErrChatNotFound,
			code:    codes.NotFound,
			message: fmt.Sprintf("chat with id %d not found", req.ChatId),
		},
		{
			name:    "internal",
			req:     req,
			expect:  message,
			service: errors.New("db is down"),
			code:    codes.Internal,
			message: "db is down",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.setup()
			s.chatService.SendMessageMock.Expect(s.ctx, tt.expect).Return(tt.service)

			resp, err := s.api.SendMessage(s.ctx, tt.req)
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
