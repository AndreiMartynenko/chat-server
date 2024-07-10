package root

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
)

func createChat(ctx context.Context, address string, certPath string, users []string) (string, error) {
	// Get access token from file
	accessToken, err := readToken()
	if err != nil {
		log.Printf("failed to read token: %v", err)
		return "", err
	}

	client, err := chatServerClient(address, certPath)
	if err != nil {
		return "", err
	}

	// Get claims from access token
	claims, err := getTokenClaims(accessToken)
	if err != nil {
		return "", err
	}

	err = isTokenExpired(claims)
	if err != nil {
		return "", err
	}

	md := metadata.New(map[string]string{"Authorization": tokenHeader + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, err := client.Create(ctx, &desc.CreateRequest{
		Chat: &desc.Chat{
			Usernames: users,
		},
	})
	if err != nil {
		return "", err
	}

	return res.GetId(), nil
}
