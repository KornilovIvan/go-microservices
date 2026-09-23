package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	accessDesc "github.com/ivankornilov/chat-server/pkg/access_v1"
)

type AccessInterceptor struct {
	client accessDesc.AccessV1Client
}

func NewAccessInterceptor(client accessDesc.AccessV1Client) *AccessInterceptor {
	return &AccessInterceptor{client: client}
}

func (i *AccessInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	_, err := i.client.Check(ctx, &accessDesc.CheckRequest{
		EndpointAddress: info.FullMethod,
	})
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
