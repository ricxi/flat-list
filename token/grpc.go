package token

import (
	"context"

	"githug.com/ricxi/flat-list/token/activation"
)

type Server struct {
	activation.UnimplementedTokenServiceServer
	R *Repository
}

func (s *Server) CreateTokenForUser(ctx context.Context, req *activation.Request) (*activation.Response, error) {
	activationToken, err := generateActivationToken()
	if err != nil {
		return nil, err
	}
	if err := s.R.InsertToken(ctx, &ActivationTokenInfo{
		Token:  activationToken,
		UserID: req.UserId,
	}); err != nil {
		return nil, err
	}

	return &activation.Response{
		ActivationToken: activationToken,
	}, nil
}
