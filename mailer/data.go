package mailer

// UserActivationData stores information needed
// to send an email to a user with instructions
// for activating their account.
type EmailActivationData struct {
	From                string `json:"from"`
	To                  string `json:"to"`
	Subject             string `json:"subject"`
	Name                string `json:"name"`
	ActivationHyperlink string `json:"activationHyperLink"`
}

type ActivationEmailData struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	ActivationData
}

type ActivationData struct {
	Name                string `json:"name"`
	ActivationHyperlink string `json:"activationHyperLink"`
}
