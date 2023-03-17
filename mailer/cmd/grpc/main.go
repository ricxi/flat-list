package main

import (
	"log"
	"net"

	"github.com/ricxi/flat-list/mailer"
	"github.com/ricxi/flat-list/mailer/pb"
	"google.golang.org/grpc"
)

func main() {
	conf, err := mailer.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":"+conf.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}

	m := mailer.NewMailer(conf.Username, conf.Password, conf.Host, conf.Port)
	mailerService := mailer.NewMailerService(m)
	srv := mailer.NewGrpcServer(mailerService)

	grpcServer := grpc.NewServer()

	pb.RegisterMailerServer(grpcServer, srv)

	log.Println("starting grpc server on port", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
