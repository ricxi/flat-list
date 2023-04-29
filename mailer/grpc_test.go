package mailer

import (
	"context"
	"testing"

	"github.com/ricxi/flat-list/mailer/pb"
	"github.com/stretchr/testify/assert"
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
