package request

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrInvalidContentType = errors.New("'Content-Type' must be 'application/json'")

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
