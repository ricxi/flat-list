// Package request contains functions for reading incoming JSON-formatted HTTP requests.
package request

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrInvalidContentType = errors.New("Invalid 'Content-Type': must be 'application/json'")

// ParseJSON decodes a JSON request body into a pointer to the given output destination.
// If the request header does not contain the header 'Content-Type: application/json', it returns an error.
func ParseJSON(r *http.Request, out any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrInvalidContentType
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&out); err != nil {
		return err
	}

	return nil
}
