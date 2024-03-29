// Package response provides functions for sending JSON-formatted HTTP responses.
package response

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Payload is a map alias type, which is used to represent JSON response data.
// Its keys must be a string type, but its values can be an any type.
type Payload map[string]any

// MustSendJSON is generally called as a last resort to send a well-formatted JSON response.
// If an error occurs while encoding the JSON payload or writing to the http.ResponseWriter,
// it will panic and set the status code to whatever was passed as an argument to the statusCode parameter.
// The caller is expected to recover from the panic, and set the appropriate status code among other things.
func MustSendJSON(w http.ResponseWriter, payload any, statusCode int, headers map[string]string) {
	if len(headers) > 0 {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	var bb bytes.Buffer
	if err := json.NewEncoder(&bb).Encode(payload); err != nil {
		panic(err)
	}

	if _, err := bb.WriteTo(w); err != nil {
		panic(err)
	}
}

// SendJSON can be called to write a JSON payload to the http.ResponseWriter.
// It can be wrapped by other functions to send more customized responses.
// If an error occurs while encoding the JSON payload or writing to the http.ResponseWriter,
// it calls SendInternalServerError to try to deliver a JSON-formatted error response.
func SendJSON(w http.ResponseWriter, payload any, statusCode int, headers map[string]string) {
	if len(headers) > 0 {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	w.Header().Set("Content-Type", "application/json")

	var bb bytes.Buffer
	if err := json.NewEncoder(&bb).Encode(payload); err != nil {
		SendInternalServerErrorAsJSON(w, err.Error())
		return
	}

	// Writing to a buffer as an extra check
	var rb bytes.Buffer
	if _, err := bb.WriteTo(&rb); err != nil {
		SendInternalServerErrorAsJSON(w, err.Error())
		return
	}

	w.WriteHeader(statusCode)
	w.Write(rb.Bytes())
}

// SendInternalServerErrorAsJSON makes one last attempt to send
// an internal server error as a JSON-formatted response. It calls
// MustSendJSON to do this, which will panic if there are any
// issues encoding or writing the message into a response.
func SendInternalServerErrorAsJSON(w http.ResponseWriter, message string) {
	msg := map[string]string{"error": message}

	MustSendJSON(w, msg, http.StatusInternalServerError, nil)
}

// SendSuccessJSON is a wrapper for SendJSON, which adds a true success flag to the JSON response body.
// If a caller attempts to add a key named "success" to the payload, it will be overwitten (not sure if I want this behaviour).
func SendSuccessJSON(w http.ResponseWriter, payload Payload, statusCode int, headers map[string]string) {
	if payload == nil {
		payload = make(Payload)
	}

	payload["success"] = true
	SendJSON(w, payload, statusCode, nil)
}

// SendErrorJSON is a wrapper for SendJSON, which adds a false success flag to the JSON response body.
func SendErrorJSON(w http.ResponseWriter, message string, statusCode int) {
	payload := Payload{
		"success": false,
		"error":   message,
	}

	SendJSON(w, payload, statusCode, nil)
}
