package config

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Metrics struct {
	ServiceName string `json:"service_name" env-default:"chat-server"`
}

type Server struct {
	Port int `json:"port" env:"SERVER_PORT" env-default:"50001"`
}
