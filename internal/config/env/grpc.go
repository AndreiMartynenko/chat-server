package env

import (
	"errors"
	"github.com/AndreiMartynenko/chat-server/internal/config"
	"os"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

// NewGRPCConfig

func GRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Get
}
