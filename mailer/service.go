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

// Service defines methods that receive email data inputs,
// and prepares and validates those inputs before calling
// methods from the Mailer type to send that data out in an email.
// It can be implemented by any infrastructure to send emails.
type Service struct {
	mailer            Mailer
	emailTemplatesDir string // this cannot be empty
}

func NewService(mailer Mailer, emailTemplatesDir string) *Service {
	return &Service{
		mailer:            mailer,
		emailTemplatesDir: emailTemplatesDir,
	}
}

// sendActivationEmail validates email data, then generates all the
// necessary templates and inputs necessary, before sending an
// activation email to a user.
// TODO: Pull validation into its own function/struct?
func (s *Service) sendActivationEmail(data ActivationEmailData) error {
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
		// this placeholder is used if no name field is provided
		data.Name = "user"
	}

	if data.Hyperlink == "" {
		return fmt.Errorf("%w: activationHyperlink", ErrMissingField)
	}

	activationEmailTmpl := filepath.Join(s.emailTemplatesDir, "useractivation.html")
	t, err := template.ParseFiles(activationEmailTmpl)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	if err := t.Execute(htmlBody, data.ActivationData); err != nil {
		return err
	}

	return s.mailer.send(data.From, data.To, data.Subject, htmlBody.String())
}
