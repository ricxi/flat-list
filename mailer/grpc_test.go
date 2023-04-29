package mailer

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/ricxi/flat-list/mailer/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// Write a mock for the Service
func TestGrpcServer_SendActivationEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := &GrpcServer{
			mailerService: &Service{
				mailer: &mockMailer{
					err: nil,
				},
				emailTemplatesDir: "./templates",
			},
		}

		req := pb.EmailRequest{
			From:    "theteam@flatlist.com",
			To:      "MichaelScott@dundermifflin.com",
			Subject: "test",
			ActivationData: &pb.ActivationData{
				Name:      "Michael",
				Hyperlink: "http://127.0.01:6000/clickme",
			},
		}

		res, err := srv.SendActivationEmail(context.Background(), &req)
		assert.NoError(t, err)

		expected := "success"
		assert.Equal(t, expected, res.Status)
	})

	t.Run("ErrorMissingField", func(t *testing.T) {
		srv := &GrpcServer{
			mailerService: &Service{
				mailer: &mockMailer{
					err: nil,
				},
				emailTemplatesDir: "./templates",
			},
		}

		req := pb.EmailRequest{
			From:    "theteam@flatlist.com",
			Subject: "test",
			ActivationData: &pb.ActivationData{
				Name:      "Michael",
				Hyperlink: "http://127.0.01:6000/clickme",
			},
		}

		res, err := srv.SendActivationEmail(context.Background(), &req)
		assert.Nil(t, res)

		expected := "missing field is required: to"
		assert.EqualError(t, err, expected)
	})
}

type cleanupFunc func(testing.TB)

// setup creates a grpc server and client
func setup(t testing.TB, service *Service) (pb.MailerClient, cleanupFunc) {
	lis := bufconn.Listen(1024 * 1024)

	srv := &GrpcServer{
		mailerService: service,
	}

	s := grpc.NewServer()
	pb.RegisterMailerServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	cc, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
	)

	if err != nil {
		t.Fatal("unable to dialto buf:", err)
	}

	return pb.NewMailerClient(cc), func(t testing.TB) {
		if err := cc.Close(); err != nil {
			t.Log("problem closing gprc client connection", err)
		}

		s.GracefulStop()
	}
}

// similar to the tests above, but tests over a mock connection
func TestGrpcServerConn_SendActivationEmail(t *testing.T) {
	t.Run("SuccessOverMockConnection", func(t *testing.T) {
		svc := Service{
			mailer: &mockMailer{
				err: nil,
			},
			emailTemplatesDir: "./templates",
		}
		c, cleanup := setup(t, &svc)
		defer cleanup(t)

		req := pb.EmailRequest{
			From:    "theteam@flatlist.com",
			To:      "MichaelScott@dundermifflin.com",
			Subject: "test",
			ActivationData: &pb.ActivationData{
				Name:      "Michael",
				Hyperlink: "http://127.0.01:6000/clickme",
			},
		}

		res, err := c.SendActivationEmail(context.Background(), &req)
		assert.NoError(t, err)

		expectedStatus := `success`
		expectedStr := `status:"success"`
		assert.Equal(t, expectedStatus, res.GetStatus())
		assert.Equal(t, expectedStr, res.String())
	})
	t.Run("InvalidArgumentErrOverMockConnection", func(t *testing.T) {
		svc := Service{
			mailer: &mockMailer{
				err: nil,
			},
			emailTemplatesDir: "./templates",
		}
		c, cleanup := setup(t, &svc)
		defer cleanup(t)

		req := pb.EmailRequest{
			From:    "theteam@flatlist.com",
			Subject: "test",
			ActivationData: &pb.ActivationData{
				Name:      "Michael",
				Hyperlink: "http://127.0.01:6000/clickme",
			},
		}

		res, err := c.SendActivationEmail(context.Background(), &req)
		assert.Nil(t, res)

		expected := `rpc error: code = InvalidArgument desc = missing field is required: to`
		assert.EqualError(t, err, expected)
	})
}
