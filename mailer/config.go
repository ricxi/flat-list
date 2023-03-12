package main

import (
  "os"
  "strconv"
)

type config struct {
	host     string
	port     int
	username string
	password string
}

func setupConfig() (*config, error) {
  conf := config{}

  conf.host = os.Getenv("HOST")
  portStr := os.Getenv("PORT")
  port, err := strconv.Atoi(portStr)
  if err != nil {
    return nil, err
  }
  conf.port = port
  conf.username = os.Getenv("USERNAME")
  conf.password = os.Getenv("PASSWORD")

  return &conf, nil
}
