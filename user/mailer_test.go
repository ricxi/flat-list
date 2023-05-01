package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	res "github.com/ricxi/flat-list/shared/response"
	"github.com/stretchr/testify/assert"
)

func Test_httpClient_sendActivationEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res.SendJSON(w, map[string]any{"success": true}, http.StatusOK, nil)
		}))
		defer ts.Close()

		mailerEndpointURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		c := httpMailerClient{
			mailerEndpointURL: *mailerEndpointURL,
		}

		err = c.sendActivationEmail(
			context.Background(),
			"michaelscott@dundermifflin.com",
			"michael",
			"activation_token_placeholder",
		)

		assert.NoError(t, err)
	})

	t.Run("ExpectedMissingFieldError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			message := "missing field is required: to"
			res.SendErrorJSON(w, message, http.StatusBadRequest)
		}))

		mailerEndpointURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		c := httpMailerClient{
			mailerEndpointURL: *mailerEndpointURL,
		}

		err = c.sendActivationEmail(
			context.Background(),
			"michaelscott@dundermifflin.com",
			"michael",
			"activation_token_placeholder",
		)

		assert.Error(t, err)
		assert.EqualError(t, err, "missing field is required: to")
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			unknown := `{"unknown":"unknown json field"}`
			w.Write([]byte(unknown))
		}))

		mailerEndpointURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		c := httpMailerClient{
			mailerEndpointURL: *mailerEndpointURL,
		}

		err = c.sendActivationEmail(
			context.Background(),
			"michaelscott@dundermifflin.com",
			"michael",
			"activation_token_placeholder",
		)

		assert.Error(t, err)
		assert.EqualError(t, err, "unknown error occurred when accessing the mailer service")
	})

	// t.Run("UnexpectedJsonDecodingError", func(t *testing.T) {
	// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		http.Error(w, "unexpected error", http.StatusInternalServerError)
	// 	}))

	// 	mailerEndpointURL, err := url.Parse(ts.URL)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	c := httpMailerClient{
	// 		mailerEndpointURL: *mailerEndpointURL,
	// 	}

	// 	err = c.sendActivationEmail(
	// 		context.Background(),
	// 		"michaelscott@dundermifflin.com",
	// 		"michael",
	// 		"activation_token_placeholder",
	// 	)

	// 	assert.Error(t, err)
	// 	assert.EqualError(t, err, "missing field is required: to")
	// })
}
