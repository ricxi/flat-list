package user

import (
	"context"

	tservice "github.com/ricxi/flat-list/token/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TokenClient interface {
	CreateActivationToken(ctx context.Context, userID string) (string, error)
	ValidateActivationToken(ctx context.Context, activationToken string) (string, error)
}

// tokenClient contains methods to call
// the token service to generate an activation token
type tokenClient struct {
	c tservice.TokenClient
}

func NewTokenClient(port string) (TokenClient, error) {
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
		return "", err
	}

	return out.UserId, nil
}
