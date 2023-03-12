package main

import (
	"bytes"
	"fmt"
	"html/template"
)

type UserActivationData struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	FirstName string `json:"firstName"`
}

type Service struct {
	tmplFilename string
	mailer       *Mailer
}

// Returns an error and does not send email if anything is missing
func (s *Service) GenerateAndSendActivationEmail(data UserActivationData) error {
	if data.From == "" {
		return fmt.Errorf("field cannot be empty: from ")
	}

	if data.To == "" {
		return fmt.Errorf("field cannot be empty: to ")
	}

	if data.FirstName == "" {
		data.FirstName = "new user"
	}

	data.Subject = "Please activate your account"

	t, err := template.ParseFiles(s.tmplFilename)
	if err != nil {
		return err
	}

	emailBody := new(bytes.Buffer)
	if err := t.Execute(emailBody, data); err != nil {
		return err
	}

	return s.mailer.Send(data.From, data.To, data.Subject, emailBody.String())
}
