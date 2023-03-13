package mailer

import (
	"context"

	"github.com/ricxi/flat-list/mailer/activate"
)

type GRPCServer struct {
	activate.UnimplementedMailerServiceServer
	EmailService *EmailService
}

// SendEmail is a grpc implementation that can be called to send an activation
// email to a user
func (s *GRPCServer) SendEmail(ctx context.Context, r *activate.Request) (*activate.Response, error) {
	data := UserActivationData{
		From:            r.From,
		To:              r.To,
		FirstName:       r.FirstName,
		ActivationToken: r.ActivationToken,
	}

	if err := s.EmailService.SendActivationEmail(data); err != nil {
		return nil, err
	}

	return &activate.Response{
		Status: "success",
	}, nil
}
