package main

import (
	"context"
	desc "github.com/AndreiMartynenko/chat-server/grpc/pkg/chat_server_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

var usernames = []string{"Bill", "Jack"}

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewChatAPIServicesV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Create(ctx, &desc.CreateNewChatRequest{Usernames: usernames})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}
	log.Printf("New chat created with ID: %d", r.Id)
	log.Printf(color.RedString("chat info: \n"), color.GreenString("%+v", r.Id))

}
