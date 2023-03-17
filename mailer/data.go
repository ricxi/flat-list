package mailer

// UserActivationData stores information needed
// to send an email to a user with instructions
// on how to activate their account
type EmailActivationData struct {
	From                string `json:"from"`
	To                  string `json:"to"`
	Name                string `json:"name"`
	ActivationHyperlink string `json:"activationHyperLink"`
}
