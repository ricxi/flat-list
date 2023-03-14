package mailer

// UserActivationData stores information needed
// to send an email a new user with instructions
// on how to activate their account
type UserActivationData struct {
	From            string `json:"from"`
	To              string `json:"to"`
	FirstName       string `json:"firstName"`
	ActivationToken string `json:"activationToken"`
}
