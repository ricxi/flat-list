package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricxi/flat-list/mailer"
	"github.com/ricxi/flat-list/shared/config"
)

// This is the starting point that can be used to compile
// the mailer into a lambda function
// ! untested
func main() {
	envs, err := config.LoadEnvs("HOST", "PORT", "USERNAME", "PASSWORD", "EMAIL_TEMPLATES", "GRPC_PORT")
	if err != nil {
		log.Fatal(err)
	}
	smtpPORT, err := strconv.Atoi(envs["PORT"])
	if err != nil {
		log.Fatal(err)
	}

	m := mailer.NewMailer(envs["USERNAME"], envs["PASSWORD"], envs["HOST"], smtpPORT)
	mailerService := mailer.NewService(m, envs["EMAIL_TEMPLATES"])

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lambda.StartWithOptions(mailer.SendActivationEmail(mailerService), lambda.WithContext(ctx))
}
