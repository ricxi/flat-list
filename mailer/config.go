package mailer

import (
	"os"
	"strconv"
)

type config struct {
	Host     string
	Port     int
	Username string
	Password string
	HttpPort string
	GrpcPort string
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
	conf.HttpPort = os.Getenv("HTTP_PORT")
	conf.GrpcPort = os.Getenv("GRPC_PORT")

	return &conf, nil
}
