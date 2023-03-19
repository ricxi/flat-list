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
	tservice "github.com/ricxi/flat-list/token/pb"
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
		return newGrpcClient(port)
	}

	return nil, errors.New("unknown client type")
}

type grpcClient struct {
	c pb.MailerClient
}

func newGrpcClient(port string) (grpcClient, error) {
	cc, err := grpc.Dial(":"+g.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return grpcClient{}, err
	}

	c := pb.NewMailerClient(cc)

	return grpcClient{
		c: c,
	}, nil
}

// SendActivationEmail makes a remote procedure call to the mailer service,
// which sends an account activation email to a newly registered user
func (g grpcClient) SendActivationEmail(email, name, activationToken string) error {
	activationHyperlink := ACTIVATION_PAGE_LINK + activationToken
	in := pb.Request{
		From:                "the.team@flat-list.com",
		To:                  email,
		Name:                name,
		ActivationHyperlink: activationHyperlink,
	}
	// maybe recompile the grpc to return a boolean?
	if _, err := g.c.SendActivationEmail(context.Background(), &in); err != nil {
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

type tokenClient struct {
	c tservice.TokenClient
}

func NewTokenClient(port string) (*tokenClient, error) {
	cc, err := grpc.Dial(":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := tservice.NewTokenClient(cc)
	return &tokenClient{
		c: c,
	}, nil
}

func (tc *tokenClient) CreateActivationToken(ctx context.Context, userID string) (string, error) {
	in := tservice.CreateTokenRequest{UserId: userID}
	out, err := tc.c.CreateActivationToken(ctx, &in)
	if err != nil {
		return "", err
	}

	return out.ActivationToken, nil
}

func (tc *tokenClient) ValidateActivationToken(ctx context.Context, activationToken string) (string, error) {
	in := tservice.ValidateTokenRequest{ActivationToken: activationToken}
	out, err := tc.c.ValidateActivationToken(context.Background(), &in)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return out.UserId, nil
}
