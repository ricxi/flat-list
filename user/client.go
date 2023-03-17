package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ricxi/flat-list/mailer"
	"github.com/ricxi/flat-list/mailer/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is used by Service to make
// http or grpc calls to other services
type Client interface {
	SendActivationEmail(email, name, activationToken string) error
}

func NewClient(clientType string) (Client, error) {
	if clientType == "http" {
		return httpClient{}, nil
	}

	if clientType == "grpc" {
		return grpcClient{}, nil
	}

	return nil, errors.New("unknown client type")
}

type grpcClient struct{}

// SendActivationEmail makes a remote procedure call to the mailer service,
// which sends an account activation email to a newly registered user
func (g grpcClient) SendActivationEmail(email, name, activationToken string) error {
	// this port should be an environment variable
	cc, err := grpc.Dial(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c := pb.NewMailerClient(cc)

	activationHyperlink := "http://localhost/9001" + activationToken
	if _, err := c.SendActivationEmail(context.Background(), &pb.Request{
		From:                "the.team@flat-list.com",
		To:                  email,
		Name:                name,
		ActivationHyperlink: activationHyperlink,
	}); err != nil {
		return err
	}

	return nil
}

type httpClient struct{}

func (h httpClient) SendActivationEmail(email, name, activationToken string) error {
	activationHyperlink := "http://localhost/9001" + activationToken
	data := mailer.EmailActivationData{
		From:                "the.team@flat-list.com",
		To:                  email,
		Name:                name,
		ActivationHyperlink: activationHyperlink,
	}

	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(&data); err != nil {
		return err
	}

	// this url should also be an environment variable or something
	req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/v1/mailer/activate", reqBody)
	if err != nil {
		return err
	}

	c := http.Client{}

	_, err = c.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
