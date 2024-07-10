package config

import (
	"github.com/joho/godotenv"
)

// GrpcConfig provides gRPC settings from config file.
type GrpcConfig interface {
	Address() string
	Transport() string
	CertPath() string
	KeyPath() string
	CaPath() string
}

// AuthConfig provides authentication service settings from config file.
type AuthConfig interface {
	Address() string
	CertPath() string
}

// TracingConfig provides tracing settings from config file.
type TracingConfig interface {
	Address() string
	ServiceName() string
}

// PgConfig provides PostgreSQL settings from config file.
type PgConfig interface {
	DSN() string
}

// Load reads .env config file.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
