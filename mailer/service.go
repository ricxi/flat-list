package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
)

var ErrMissingField = errors.New("missing field is required")

const (
	// Subject line for the user activation email
	activationEmailSubject string = "Please activate your account"
)

// MailerService defines methods that receive email data inputs,
// and prepares and validates those inputs before calling
// methods from the Mailer type to send that data out in an email.
// It can be implemented by any infrastructure to send emails.
type MailerService struct {
	mailer            Mailer
	emailTemplatesDir string
}

func NewMailerService(mailer Mailer, emailTemplatesDir string) *MailerService {
	return &MailerService{
		mailer:            mailer,
		emailTemplatesDir: emailTemplatesDir,
	}
}

// SendActivationEmail validates email data, then generates all the
// necessary templates and inputs necessary, before sending an
// activation email to a user.
// TODO: Pull validation into its own function/struct?
func (s *MailerService) sendActivationEmail(data ActivationEmailData) error {
	if data.From == "" {
		return fmt.Errorf("%w: from", ErrMissingField)
	}

	if data.To == "" {
		return fmt.Errorf("%w: to", ErrMissingField)
	}

	if data.Subject == "" {
		return fmt.Errorf("%w: subject", ErrMissingField)
	}

	if data.Name == "" {
		// a placeholder is used if no name field is provided
		data.Name = "user"
	}

	if data.ActivationHyperlink == "" {
		return fmt.Errorf("%w: hyperlink", ErrMissingField)
	}

	activationEmailTmpl := filepath.Join(s.emailTemplatesDir, "useractivation.html")
	t, err := template.ParseFiles(activationEmailTmpl)
	if err != nil {
		return err
	}

	// tmplData := struct {
	// 	Name                string
	// 	ActivationHyperlink string
	// }{
	// 	Name:                data.Name,
	// 	ActivationHyperlink: data.ActivationHyperlink,
	// }

	htmlBody := new(bytes.Buffer)
	if err := t.Execute(htmlBody, data.ActivationData); err != nil {
		return err
	}

	return s.mailer.send(data.From, data.To, data.Subject, htmlBody.String())
}
