package auth

import (
	"context"

	desc "github.com/AndreiMartynenko/auth/pkg/access_v1"
	"github.com/AndreiMartynenko/chat-server/internal/client/rpc"
)

type authClient struct {
	client desc.AccessV1Client
}

var _ rpc.AuthClient = (*authClient)(nil)

// NewAuthClient creates new AuthClient object.
func NewAuthClient(client desc.AccessV1Client) rpc.AuthClient {
	return &authClient{
		client: client,
	}
}

// Check calls authentication service method for authorization.
func (c *authClient) Check(ctx context.Context, endpoint string) error {
	_, err := c.client.Check(ctx, &desc.CheckRequest{
		Endpoint: endpoint,
	})
	return err
}
