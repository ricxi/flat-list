package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/ricxi/flat-list/mailer"
	"github.com/ricxi/flat-list/mailer/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ActivationPageLink string = "http://localhost:5173/activate?token="

// MailerClient is used by Service to make
// http or grpc calls to other services
type MailerClient interface {
	sendActivationEmail(ctx context.Context, email, name, activationToken string) error
}

type grpcMailerClient struct {
	c pb.MailerClient
}

func NewGRPCMailerClient(port string) (*grpcMailerClient, error) {
	cc, err := grpc.Dial(":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewMailerClient(cc)

	return &grpcMailerClient{
		c: c,
	}, nil
}

// SendActivationEmail makes a remote procedure call to the mailer service,
// which sends an account activation email to a newly registered user
func (g *grpcMailerClient) sendActivationEmail(ctx context.Context, email, name, activationToken string) error {
	activationHyperlink := ActivationPageLink + activationToken
	in := pb.EmailRequest{
		From:    "the.team@flat-list.com",
		To:      email,
		Subject: "Please activate your account",
		ActivationData: &pb.ActivationData{
			Name:      name,
			Hyperlink: activationHyperlink,
		},
	}
	if _, err := g.c.SendActivationEmail(ctx, &in); err != nil {
		return err
	}

	return nil
}

type httpMailerClient struct {
	mailerEndpointURL url.URL
}

func NewHTTPMailerClient(mailerEndpoint string) (*httpMailerClient, error) {
	mailerEndpointURL, err := url.Parse(mailerEndpoint)
	if err != nil {
		return nil, err
	}

	return &httpMailerClient{
		mailerEndpointURL: *mailerEndpointURL,
	}, nil
}

func (h *httpMailerClient) sendActivationEmail(ctx context.Context, email, name, activationToken string) error {
	activationHyperlink := ActivationPageLink + activationToken

	data := mailer.ActivationEmailData{
		From:    "the.team@flat-list.com",
		To:      email,
		Subject: "Follow the instructions to activate your account",
		ActivationData: mailer.ActivationData{
			Name:      name,
			Hyperlink: activationHyperlink,
		},
	}

	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(&data); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.mailerEndpointURL.String(), reqBody)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	c := http.Client{Timeout: 5 * time.Second}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		// Custom utility to extract errors?
		errs := struct {
			ErrStr string `json:"error"`
		}{}

		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&errs); err != nil {
			return err
		}

		if errs.ErrStr != "" {
			return errors.New(errs.ErrStr)
		}

		return errors.New("unknown error occurred when accessing the mailer service")
	}

	return nil
}
