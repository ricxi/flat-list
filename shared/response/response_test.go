package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ricxi/flat-list/shared/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// superfluousHeaderDetector is a wrapper for the
// httptest.ResponseRecorder type which is used to
// count the number of times WriteHeader was called
type superfluousHeaderDetector struct {
	*httptest.ResponseRecorder
	count int
}

func (r *superfluousHeaderDetector) WriteHeader(statusCode int) {
	r.count++
	r.ResponseRecorder.WriteHeader(statusCode)
}

func TestMustSendJSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)
		rr := httptest.NewRecorder()

		expectedResponse := `{"id":"abcdefghijkl","success":true}`
		payload := map[string]any{
			"success": true,
			"id":      "abcdefghijkl",
		}

		err := response.MustSendJSON(rr, payload, http.StatusOK, nil)
		require.NoError(err)

		assert.Equal(http.StatusOK, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expectedResponse, rr.Body.String())
		}
	})

	t.Run("FailSuperfluousResponseWriteHeaderCall", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		rr := &superfluousHeaderDetector{ResponseRecorder: httptest.NewRecorder()}

		expectedHeaderCalls := 1

		invalidPayload := make(chan string)

		err := response.MustSendJSON(rr, invalidPayload, http.StatusOK, nil)
		require.Error(err)

		assert.Equal(expectedHeaderCalls, rr.count, "got more than one WriteHeader calls, but expected one")
	})
}

func TestErrorJSON(t *testing.T) {
	t.Run("FailSuperfluousResponseWriteHeaderCall", func(t *testing.T) {
		assert := assert.New(t)

		rr := &superfluousHeaderDetector{ResponseRecorder: httptest.NewRecorder()}

		expectedHeaderCalls := 1

		invalidPayload := make(chan string)

		response.ErrorJSON(rr, http.StatusOK, invalidPayload)

		assert.Equal(expectedHeaderCalls, rr.count, "got more than one WriteHeader calls, but expected one")
	})
}
