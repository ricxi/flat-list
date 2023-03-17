package mailer

import (
	"gopkg.in/gomail.v2"
)

// Mailer defines methods for receiving data and sending emails.
type Mailer struct {
	dialer *gomail.Dialer
}

func NewMailer(username, password, host string, port int) *Mailer {
	dialer := gomail.NewDialer(host, port, username, password)

	return &Mailer{
		dialer: dialer,
	}
}

// Send is a wrapper for SendMultiple.
// It sends an email to a single recipient
func (m *Mailer) Send(from, to, subject, body string) error {
	return m.SendMultiple(from, subject, body, to)
}

// SendMultiple sends an email to one or more recipients
func (m *Mailer) SendMultiple(from, subject, body string, recipients ...string) error {
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
