package mailer

import (
	"context"

	"github.com/ricxi/flat-list/mailer/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedMailerServer
	mailerService *Service
}

func NewGrpcServer(mailerService *Service) GrpcServer {
	return GrpcServer{
		mailerService: mailerService,
	}
}

// SendActivationEmail is a grpc implementation that can be called by other
// services to send an activation email to a user.
func (gs GrpcServer) SendActivationEmail(ctx context.Context, r *pb.EmailRequest) (*pb.Response, error) {
	data := ActivationEmailData{
		From:    r.From,
		To:      r.To,
		Subject: activationEmailSubject,
		ActivationData: ActivationData{
			Name:      r.ActivationData.Name,
			Hyperlink: r.ActivationData.Hyperlink,
		},
	}
	if err := gs.mailerService.sendActivationEmail(data); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.Response{
		Status: "success",
	}, nil
}
