package main

// UserActivationData stores information needed
// to send an activation email to validate a new user
type UserActivationData struct {
	From            string `json:"from"`
	To              string `json:"to"`
	FirstName       string `json:"firstName"`
	ActivationToken string `json:"activationToken"`
}
