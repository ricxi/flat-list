package mailer

import (
	"context"

	"github.com/ricxi/flat-list/mailer/pb"
)

type GrpcServer struct {
	pb.UnimplementedMailerServer
	EmailService *EmailService
}

// SendActivationEmail is a grpc implementation that can be called by other
// services to send an activation email to a user
func (gs GrpcServer) SendActivationEmail(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	data := UserActivationData{
		From:            r.From,
		To:              r.To,
		FirstName:       r.Name,
		ActivationToken: r.ActivationToken,
	}

	if err := gs.EmailService.SendActivationEmail(data); err != nil {
		return nil, err
	}

	return &pb.Response{
		Status: "success",
	}, nil
}
