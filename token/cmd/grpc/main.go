package main

import (
	"log"
	"net"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"githug.com/ricxi/flat-list/token"
	"githug.com/ricxi/flat-list/token/activation"
	"google.golang.org/grpc"
)

func main() {
	GrpcPort := os.Getenv("GRPC_PORT")
	if GrpcPort == "" {
		log.Fatalln("missing grpc env")
	}

	config, err := token.LoadConfig()
	if err != nil {
		log.Fatalln("problem loading configuation: ", err)
	}

	db, err := token.Connect(config.DatabaseURL)
	if err != nil {
		log.Fatalln("problem connecting to postgres: ", err)
	}
	defer db.Close()

	repo := token.Repository{
		DB: db,
	}

	lis, err := net.Listen("tcp", ":"+GrpcPort)
	if err != nil {
		log.Fatalln("fail to listen on tcp", err)
	}

	grpcServer := grpc.NewServer()
	srv := token.Server{R: &repo}

	activation.RegisterTokenServiceServer(grpcServer, &srv)

	log.Println("starting grpc token server on ", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("failed to start grpc server ", err)
	}
}
