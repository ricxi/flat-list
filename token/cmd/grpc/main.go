package main

import (
	"log"
	"net"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ricxi/flat-list/shared/config"
	"github.com/ricxi/flat-list/token"
	"github.com/ricxi/flat-list/token/pb"
	"google.golang.org/grpc"
)

func main() {
	envs, err := config.LoadEnvs("DATABASE_URL", "GRPC_PORT")
	if err != nil {
		log.Fatalln("problem loading configuration: ", err)
	}

	db, err := token.Connect(envs["DATABASE_URL"])
	if err != nil {
		log.Fatalln("problem connecting to postgres: ", err)
	}
	defer db.Close()

	repo := token.NewRepository(db)

	lis, err := net.Listen("tcp", ":"+envs["GRPC_PORT"])
	if err != nil {
		log.Fatalln("fail to listen on tcp", err)
	}

	grpcServer := grpc.NewServer()
	srv := token.Server{Repository: repo}

	pb.RegisterTokenServer(grpcServer, srv)

	log.Println("starting grpc token server on", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("failed to start grpc server ", err)
	}
}
