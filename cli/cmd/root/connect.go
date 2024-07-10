package root

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
)

func connectChat(ctx context.Context, address string, certPath string, chatID string) error {
	// Get access token from file
	accessToken, err := readToken()
	if err != nil {
		log.Printf("failed to read token: %v", err)
		return err
	}

	client, err := chatServerClient(address, certPath)
	if err != nil {
		return err
	}

	// Get claims from access token
	claims, err := getTokenClaims(accessToken)
	if err != nil {
		return err
	}

	err = isTokenExpired(claims)
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"Authorization": tokenHeader + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := client.Connect(ctx, &desc.ConnectRequest{
		ChatId:   chatID,
		Username: claims.Username,
	})
	if err != nil {
		return err
	}

	fmt.Println(color.CyanString("[Connected]"))

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Printf("failed to receive message: %v", errRecv)
				return
			}

			fmt.Printf("[%v] - [from: %s]: %s",
				color.YellowString(message.GetTimestamp().AsTime().Format(time.RFC3339)),
				color.GreenString(message.GetFrom()),
				message.GetText(),
			)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		// Input will be replaced with formatted message when received
		fmt.Print(lineUp)
		_, err = client.SendMessage(ctx, &desc.SendMessageRequest{
			ChatId: chatID,
			Message: &desc.Message{
				From:      claims.Username,
				Text:      text,
				Timestamp: timestamppb.Now(),
			},
		})
		if err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}
}
