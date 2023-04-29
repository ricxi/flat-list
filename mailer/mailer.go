package mailer

import (
	"gopkg.in/gomail.v2"
)

// Mailer defines methods for receiving data to send emails.
type Mailer interface {
	send(from, to, subject, body string) error
}

// mailer implements the gomail package to send emails
type mailer struct {
	dialer *gomail.Dialer
}

func NewMailer(username, password, host string, port int) Mailer {
	dialer := gomail.NewDialer(host, port, username, password)

	return &mailer{
		dialer: dialer,
	}
}

// send is a wrapper for SendMultiple.
// It sends an email to a single recipient
func (m *mailer) send(from, to, subject, body string) error {
	return m.sendMultiple(from, subject, body, to)
}

// sendMultiple sends an email to one or more recipients
func (m *mailer) sendMultiple(from, subject, body string, recipients ...string) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", from)
	msg.SetHeader("To", recipients...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
