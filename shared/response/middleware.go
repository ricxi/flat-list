package response

import "net/http"

const defaultMessage string = `{"error": "an internal server error has occurred"}`

// HandlePanic recovers from a panic and sends a default
// error message with a 500 status code. It is intended to
// be used with MustSendAsJSON.
func HandlePanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(defaultMessage))
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
