package env

import (
	"errors"
	"github.com/AndreiMartynenko/chat-server/internal/config"
	"os"
)

var _ config.PGConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

//NewPGConfig

func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
