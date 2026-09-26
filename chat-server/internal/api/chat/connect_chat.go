package chat

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

func (i *Implementation) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer) error {
	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetChatId()]
	i.mxChannel.RUnlock()
	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	i.mxChat.Lock()
	if _, okChat := i.chats[req.GetChatId()]; !okChat {
		i.chats[req.GetChatId()] = &chatStreams{
			streams: make(map[string]desc.ChatV1_ConnectChatServer),
		}
	}
	chat := i.chats[req.GetChatId()]
	i.mxChat.Unlock()

	chat.m.Lock()
	chat.streams[req.GetUsername()] = stream
	chat.m.Unlock()

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}

			chat.m.RLock()
			for _, st := range chat.streams {
				if err := st.Send(msg); err != nil {
					chat.m.RUnlock()
					return err
				}
			}
			chat.m.RUnlock()
		case <-stream.Context().Done():
			chat.m.Lock()
			delete(chat.streams, req.GetUsername())
			chat.m.Unlock()
			return nil
		}
	}
}
