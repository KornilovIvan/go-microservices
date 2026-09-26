package chat

import (
	"sync"

	"github.com/ivankornilov/chat-server/internal/service"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

type chatStreams struct {
	streams map[string]desc.ChatV1_ConnectChatServer
	m       sync.RWMutex
}

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService

	chats  map[int64]*chatStreams
	mxChat sync.RWMutex

	channels  map[int64]chan *desc.Message
	mxChannel sync.RWMutex
}

func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
		chats:       make(map[int64]*chatStreams),
		channels:    make(map[int64]chan *desc.Message),
	}
}
