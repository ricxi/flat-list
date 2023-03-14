package main

import (
	"log"
	"net"

	"github.com/ricxi/flat-list/mailer"
	"github.com/ricxi/flat-list/mailer/activate"
	"google.golang.org/grpc"
)

func main() {
	conf, err := mailer.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	m := mailer.NewMailer(conf.Username, conf.Password, conf.Host, conf.Port)
	es := mailer.NewEmailService(m)
	srv := mailer.GRPCServer{EmailService: es}

	lis, err := net.Listen("tcp", ":"+conf.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()

	activate.RegisterMailerServiceServer(grpcServer, &srv)

	log.Println("starting grpc server on port", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
