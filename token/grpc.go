package token

import (
	"context"

	"github.com/ricxi/flat-list/token/pb"
)

type Server struct {
	pb.UnimplementedTokenServer
	Repository
}

func (s Server) CreateActivationToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.CreateTokenResponse, error) {
	activationToken, err := generateActivationToken()
	if err != nil {
		return nil, err
	}
	if err := s.insertActivationToken(ctx, &ActivationTokenInfo{
		Token:  activationToken,
		UserID: req.UserId,
	}); err != nil {
		return nil, err
	}

	return &pb.CreateTokenResponse{
		ActivationToken: activationToken,
	}, nil
}

func (s Server) ValidateActivationToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, err := s.getUserID(ctx, req.ActivationToken)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateTokenResponse{
		UserId: userID,
	}, nil
}
