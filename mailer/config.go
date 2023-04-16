package mailer

import (
	"os"
	"strconv"
)

type smtp struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
}

type envs struct {
	Smtp              smtp   `mapstructure:"SMTP"`
	EmailTemplatesDir string `mapstructure:"EMAIL_TEMPLATES"`
	HttpPort          string `mapstructure:"HTTP_PORT"`
	GrpcPort          string `mapstructure:"GRPC_PORT"`
}

type config struct {
	Host              string
	Port              int
	Username          string
	Password          string
	EmailTemplatesDir string
	HttpPort          string
	GrpcPort          string
}

func SetupConfig() (*config, error) {
	conf := config{}

	conf.Host = os.Getenv("HOST")

	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	conf.Port = port

	conf.Username = os.Getenv("USERNAME")
	conf.Password = os.Getenv("PASSWORD")
	conf.EmailTemplatesDir = os.Getenv("EMAIL_TEMPLATES")

	conf.HttpPort = os.Getenv("HTTP_PORT")
	conf.GrpcPort = os.Getenv("GRPC_PORT")

	return &conf, nil
}
