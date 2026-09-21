package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/ivankornilov/auth/internal/config"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	flag.Parse()

	cfg, err := config.Parse(configPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	conn, err := grpc.Dial(cfg.GRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createRes, err := c.Create(ctx, &desc.CreateRequest{
		Name:            "Ivan",
		Email:           "ivan@example.com",
		Password:        "qwerty",
		PasswordConfirm: "qwerty",
		Role:            desc.Role_USER,
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Println(color.GreenString("Create response: %+v", createRes))

	getRes, err := c.Get(ctx, &desc.GetRequest{Id: createRes.GetId()})
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	log.Println(color.GreenString("Get response: %+v", getRes))

	_, err = c.Update(ctx, &desc.UpdateRequest{
		Id:    createRes.GetId(),
		Name:  wrapperspb.String("Ivan Updated"),
		Email: wrapperspb.String("ivan.updated@example.com"),
	})
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}
	log.Println(color.GreenString("Update: ok"))

	getRes, err = c.Get(ctx, &desc.GetRequest{Id: createRes.GetId()})
	if err != nil {
		log.Fatalf("failed to get updated user: %v", err)
	}
	log.Println(color.GreenString("Get after update: %+v", getRes))

	_, err = c.Delete(ctx, &desc.DeleteRequest{Id: createRes.GetId()})
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}
	log.Println(color.GreenString("Delete: ok"))
}
