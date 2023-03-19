package token

import (
	"context"

	"github.com/ricxi/flat-list/token/pb"
)

type Server struct {
	pb.UnimplementedTokenServer
	Repository
}

func (s Server) CreateActivationToken(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	activationToken, err := generateActivationToken()
	if err != nil {
		return nil, err
	}
	if err := s.InsertActivationToken(ctx, &ActivationTokenInfo{
		Token:  activationToken,
		UserID: req.UserId,
	}); err != nil {
		return nil, err
	}

	return &pb.Response{
		ActivationToken: activationToken,
	}, nil
}
