package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var ErrNilPayload = errors.New("payload cannot be nil")

// MustSendJSON is a helper function used by all methods in the response package to write a JSON payload
// to the http.ResponseWriter. If an error occurs while encoding the JSON data or writing to the
// http.ResponseWriter, it will panic with the purpose of being recovered from by the caller so that
// the correct status code can be set.
func MustSendJSON(w http.ResponseWriter, payload any, statusCode int, headers map[string]string) error {
	if payload == nil {
		// should I panic instead? I feel like the caller should not be able to pass a nil payload
		return ErrNilPayload
	}

	if len(headers) > 0 {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	w.WriteHeader(statusCode)

	var bb bytes.Buffer
	err := json.NewEncoder(&bb).Encode(payload)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	if _, err := bb.WriteTo(w); err != nil {
		return err
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, status int, message any) {
	env := map[string]any{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := MustSendJSON(w, env, status, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}
