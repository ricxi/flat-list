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

const ACTIVATION_PAGE_LINK string = "http://localhost:5173/activate?token="

// Client is used by Service to make
// http or grpc calls to other services
type Client interface {
	SendActivationEmail(email, name, activationToken string) error
}

// NewMailerClient can be called to create a grpc or http mailer client
func NewMailerClient(clientType, port string) (Client, error) {
	if clientType == "http" {
		return httpClient{}, nil
	}

	if clientType == "grpc" {
		return grpcClient{port: port}, nil
	}

	return nil, errors.New("unknown client type")
}

type grpcClient struct {
	port string
}

// SendActivationEmail makes a remote procedure call to the mailer service,
// which sends an account activation email to a newly registered user
func (g grpcClient) SendActivationEmail(email, name, activationToken string) error {
	// this port should be an environment variable
	cc, err := grpc.Dial(":"+g.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c := pb.NewMailerClient(cc)

	activationHyperlink := ACTIVATION_PAGE_LINK + activationToken
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

type httpClient struct {
	port string
}

func (h httpClient) SendActivationEmail(email, name, activationToken string) error {
	activationHyperlink := ACTIVATION_PAGE_LINK + activationToken
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

	// this is kind of sketchy right now, but I'll fix it later
	req, err := http.NewRequest(http.MethodPost, "http://localhost:"+h.port+"/v1/mailer/activate", reqBody)
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
