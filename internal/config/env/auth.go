package env

import (
	"errors"
	"net"
	"os"

	"github.com/AndreiMartynenko/chat-server/internal/config"
)

const (
	authHostEnvName     = "AUTH_HOST"
	authPortEnvName     = "AUTH_PORT"
	authCertPathEnvName = "AUTH_CERT_PATH"
)

type authConfig struct {
	host     string
	port     string
	certPath string
}

var _ config.AuthConfig = (*authConfig)(nil)

// NewAuthConfig creates new object of AuthConfig interface.
func NewAuthConfig() (config.AuthConfig, error) {
	host := os.Getenv(authHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("authentication service host not found")
	}

	port := os.Getenv(authPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("authentication service port not found")
	}

	certPath := os.Getenv(authCertPathEnvName)
	if len(certPath) == 0 {
		return nil, errors.New("authentication service certificate not found")
	}

	return &authConfig{
		host:     host,
		port:     port,
		certPath: certPath,
	}, nil
}

func (cfg *authConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *authConfig) CertPath() string {
	return cfg.certPath
}
