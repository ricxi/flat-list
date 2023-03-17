package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricxi/flat-list/mailer"
)

// This is the starting point that can be used to compile
// the mailer into a lambda function
func main() {
	conf, err := mailer.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	m := mailer.NewMailer(conf.Username, conf.Password, conf.Host, conf.Port)
	es := mailer.NewEmailService(m)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lambda.StartWithOptions(mailer.SendActivationEmail(es), lambda.WithContext(ctx))
}
