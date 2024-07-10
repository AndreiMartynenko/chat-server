package rpc

import (
	"context"
)

// AuthClient is a client for authentication service.
type AuthClient interface {
	Check(ctx context.Context, endpoint string) error
}
