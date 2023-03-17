package mailer

import (
	"bytes"
	"fmt"
	"html/template"
)

const (
	// Subject line for the user activation email
	ACTIVATION_EMAIL_SUBJECT string = "Please activate your account"
	// Path to the HTML template used to generate the body of the user activation email
	ACTIVATION_HTML_TMPL string = "./templates/useractivation.html"
)

// MailerService defines methods that receive email data inputs,
// and prepares and validates those inputs before calling
// methods from the Mailer type to send that data out in an email.
// It can be implemented by any infrastructure to send emails.
type MailerService struct {
	mailer *Mailer
}

func NewMailerService(mailer *Mailer) *MailerService {
	return &MailerService{
		mailer: mailer,
	}
}

// SendActivationEmail validates email data, then generates all the
// necessary templates and inputs necessary, before sending an
// activation email to a user.
func (s *MailerService) SendActivationEmail(data EmailActivationData) error {
	if data.From == "" {
		return fmt.Errorf("field cannot be empty: from ")
	}

	if data.To == "" {
		return fmt.Errorf("field cannot be empty: to ")
	}

	if data.Name == "" {
		data.Name = "user"
	}

	if data.ActivationHyperlink == "" {
		return fmt.Errorf("field cannot be empty: hyperlink ")
	}

	t, err := template.ParseFiles(ACTIVATION_HTML_TMPL)
	if err != nil {
		return err
	}

	tmplData := struct {
		Name                string
		ActivationHyperlink string
	}{
		Name:                data.Name,
		ActivationHyperlink: data.ActivationHyperlink,
	}

	htmlBody := new(bytes.Buffer)
	if err := t.Execute(htmlBody, tmplData); err != nil {
		return err
	}

	return s.mailer.Send(data.From, data.To, ACTIVATION_EMAIL_SUBJECT, htmlBody.String())
}
