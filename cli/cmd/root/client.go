package root

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	descAuth "github.com/AndreiMartynenko/auth/pkg/auth_v1"
	descChat "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
	"github.com/AndreiMartynenko/common/pkg/closer"
)

func authClient(address string, certPath string) (descAuth.AuthV1Client, error) {
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, err
	}
	closer.Add(conn.Close)

	return descAuth.NewAuthV1Client(conn), nil
}

func chatServerClient(address string, certPath string) (descChat.ChatV1Client, error) {
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, err
	}
	closer.Add(conn.Close)

	return descChat.NewChatV1Client(conn), nil
}
