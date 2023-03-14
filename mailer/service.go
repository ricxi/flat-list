package mailer

import (
	"bytes"
	"fmt"
	"html/template"
)

// Subject line for user activation emails
const (
	ACTIVATION_EMAIL_SUBJECT  string = "Please activate your account"
	ACTIVATION_HTML_TMPL_NAME string = "./useractivation.html"
)

// EmailService for sending emails
type EmailService struct {
	mailer *Mailer
}

func NewEmailService(mailer *Mailer) *EmailService {
	return &EmailService{
		mailer: mailer,
	}
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
		data.FirstName = "user"
	}

	t, err := template.ParseFiles(ACTIVATION_HTML_TMPL_NAME)
	if err != nil {
		return err
	}

	activationLink := "http://localhost:5000/" + data.ActivationToken

	emailTemplateData := struct {
		FirstName      string `json:"firstName"`
		ActivationLink string `json:"activationLink"`
	}{
		FirstName:      data.FirstName,
		ActivationLink: activationLink,
	}

	htmlEmailBody := new(bytes.Buffer)
	if err := t.Execute(htmlEmailBody, emailTemplateData); err != nil {
		return err
	}

	subject := ACTIVATION_EMAIL_SUBJECT
	return s.mailer.Send(data.From, data.To, subject, htmlEmailBody.String())
}
