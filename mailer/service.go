package main

import (
	"bytes"
	"fmt"
	"html/template"
)

// Subject line for user activation emails
const ACTIVATION_EMAIL_SUBJECT string = "Please activate your account"

// EmailService for sending emails
type EmailService struct {
	tmplFilename string
	mailer       *Mailer
}

// SendActivationEmail validates that all the data is provided to send
// an activation email, generates an html template for the email, and
// calls the mailer to send the email
func (s *EmailService) SendActivationEmail(data UserActivationData) error {
	if data.From == "" {
		return fmt.Errorf("field cannot be empty: from ")
	}

	if data.To == "" {
		return fmt.Errorf("field cannot be empty: to ")
	}

	if data.FirstName == "" {
		data.FirstName = "new user"
	}

	t, err := template.ParseFiles(s.tmplFilename)
	if err != nil {
		return err
	}

	activationToken, err := generateActivationToken()
	if err != nil {
		return err
	}

	data.Subject = ACTIVATION_EMAIL_SUBJECT
	data.Token = activationToken

	htmlEmailBody := new(bytes.Buffer)
	if err := t.Execute(htmlEmailBody, data); err != nil {
		return err
	}

	return s.mailer.Send(data.From, data.To, data.Subject, htmlEmailBody.String())
}
