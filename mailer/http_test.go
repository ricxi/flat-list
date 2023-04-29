package mailer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleSendActivationEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service := Service{
			mailer: &mockMailer{
				err: nil,
			},
			emailTemplatesDir: "./templates",
		}
		h := HandleSendActivationEmail(&service)

		rr := httptest.NewRecorder()

		body := `
		{
			"from":    "theteam@flatlist.com",
			"to":      "MichaelScott@dundermifflin.com",
			"subject": "test",
			"activationData": {
				"name":      "Michael",
				"hyperlink": "http://127.0.01:6000/clickme"
			}
		},`

		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

		h.ServeHTTP(rr, r)

		assert.Equal(t, http.StatusOK, rr.Code)

		expected := `{"success": true}`
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("ErrorResponse", func(t *testing.T) {
		service := Service{
			mailer: &mockMailer{
				err: nil,
			},
			emailTemplatesDir: "./templates",
		}
		h := HandleSendActivationEmail(&service)

		rr := httptest.NewRecorder()

		body := `
		{
			"from":    "theteam@flatlist.com",
			"subject": "test",
			"activationData": {
				"name":      "Michael",
				"hyperlink": "http://127.0.01:6000/clickme"
			}
		},`

		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

		h.ServeHTTP(rr, r)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		expected := `{"error":"missing field is required: to"}`
		assert.JSONEq(t, expected, rr.Body.String())
	})
}
