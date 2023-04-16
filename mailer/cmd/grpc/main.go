package main

import (
	"log"
	"net"
	"strconv"

	"github.com/ricxi/flat-list/mailer"
	"github.com/ricxi/flat-list/mailer/pb"
	"github.com/ricxi/flat-list/shared/config"
	"google.golang.org/grpc"
)

func main() {
	envs, err := config.LoadEnvs("HOST", "PORT", "USERNAME", "PASSWORD", "EMAIL_TEMPLATES", "GRPC_PORT")
	if err != nil {
		log.Fatal(err)
	}
	smtpPORT, err := strconv.Atoi(envs["PORT"])
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":"+envs["GRPC_PORT"])
	if err != nil {
		log.Fatal(err)
	}

	m := mailer.NewMailer(envs["USERNAME"], envs["PASSWORD"], envs["HOST"], smtpPORT)
	mailerService := mailer.NewMailerService(m, envs["EMAIL_TEMPLATES"])
	srv := mailer.NewGrpcServer(mailerService)

	grpcServer := grpc.NewServer()

	pb.RegisterMailerServer(grpcServer, srv)

	log.Println("starting grpc mailer server for on port", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
