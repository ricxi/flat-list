package mailer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMailer struct {
	err error
}

func (m *mockMailer) send(from, to, subject, body string) error {
	fmt.Println(body)
	return m.err
}

func TestServiceSendActivationEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mailer := mockMailer{err: nil}
		service := &MailerService{
			mailer:            &mailer,
			emailTemplatesDir: "./templates",
		}

		data := ActivationEmailData{
			To:      "michaelscott@dundermifflin.com",
			From:    "theteam@flatlist.com",
			Subject: "How to activate your new account",
			ActivationData: ActivationData{
				Name:                "Michael",
				ActivationHyperlink: "http://localhost:5000/clickme",
			},
		}

		err := service.sendActivationEmail(data)

		assert.NoError(t, err)
	})

}
