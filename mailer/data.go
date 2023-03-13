package main

// UserActivationData stores information needed
// to send an activation email to validate a new user
type UserActivationData struct {
	From           string `json:"from"`
	To             string `json:"to"`
	Subject        string `json:"subject"`
	FirstName      string `json:"firstName"`
	ActivationLink string `json:"activationLink"`
}
