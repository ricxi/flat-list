package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
)

const (
	// Subject line for the user activation email
	activationEmailSubject string = "Please activate your account"
)

// MailerService defines methods that receive email data inputs,
// and prepares and validates those inputs before calling
// methods from the Mailer type to send that data out in an email.
// It can be implemented by any infrastructure to send emails.
type MailerService struct {
	mailer            *Mailer
	emailTemplatesDir string
}

func NewMailerService(mailer *Mailer, emailTemplatesDir string) *MailerService {
	return &MailerService{
		mailer:            mailer,
		emailTemplatesDir: emailTemplatesDir,
	}
}

// SendActivationEmail validates email data, then generates all the
// necessary templates and inputs necessary, before sending an
// activation email to a user.
// TODO: Add validation
func (s *MailerService) sendActivationEmail(data EmailActivationData) error {
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

	activationEmailTmpl := filepath.Join(s.emailTemplatesDir, "useractivation.html")
	t, err := template.ParseFiles(activationEmailTmpl)
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

	return s.mailer.send(data.From, data.To, activationEmailSubject, htmlBody.String())
}
