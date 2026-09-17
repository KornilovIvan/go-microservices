package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ivankornilov/chat-server/internal/config"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	conn, err := grpc.Dial(grpcConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewChatV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createRes, err := c.Create(ctx, &desc.CreateRequest{
		Usernames: []string{"ivan", "anna"},
	})
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}
	log.Printf(color.GreenString("Create response: %+v", createRes))

	_, err = c.SendMessage(ctx, &desc.SendMessageRequest{
		ChatId:    createRes.GetId(),
		From:      "ivan",
		Text:      "hello",
		Timestamp: timestamppb.Now(),
	})
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}
	log.Printf(color.GreenString("SendMessage: ok"))

	_, err = c.Delete(ctx, &desc.DeleteRequest{Id: createRes.GetId()})
	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}
	log.Printf(color.GreenString("Delete: ok"))
}
