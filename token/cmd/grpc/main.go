package main

import (
	"log"
	"net"

	_ "github.com/jackc/pgx/v5/stdlib"
	"githug.com/ricxi/flat-list/token"
	"githug.com/ricxi/flat-list/token/pb"
	"google.golang.org/grpc"
)

func main() {
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

	lis, err := net.Listen("tcp", ":"+config.GrpcPort)
	if err != nil {
		log.Fatalln("fail to listen on tcp", err)
	}

	grpcServer := grpc.NewServer()
	srv := token.Server{Repository: repo}

	pb.RegisterTokenServer(grpcServer, srv)

	log.Println("starting grpc token server on ", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("failed to start grpc server ", err)
	}
}
