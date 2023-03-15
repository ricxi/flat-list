package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ricxi/flat-list/mailer/activate"
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
	cc, err := grpc.Dial(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c := activate.NewMailerServiceClient(cc)
	if _, err := c.SendEmail(context.Background(), &activate.Request{
		From:            "the.team@flat-list.com",
		To:              email,
		FirstName:       name,
		ActivationToken: activationToken,
	}); err != nil {
		return err
	}

	return nil
}

type httpClient struct{}

func (h httpClient) SendActivationEmail(email, firstName, activationToken string) error {
	activationInfo := struct {
		From            string `json:"from"`
		To              string `json:"to"`
		FirstName       string `json:"firstName"`
		ActivationToken string `json:"activationToken"`
	}{
		From:            "the.team@flatlist.com",
		To:              email,
		FirstName:       firstName,
		ActivationToken: activationToken,
	}

	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(&activationInfo); err != nil {
		return err
	}

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
