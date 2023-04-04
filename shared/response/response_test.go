package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ricxi/flat-list/shared/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeHeaderRecorder is a wrapper for the
// httptest.ResponseRecorder type which is used to
// count the number of times WriteHeader was called
type writeHeaderRecorder struct {
	*httptest.ResponseRecorder
	Count int
}

func (r *writeHeaderRecorder) WriteHeader(statusCode int) {
	r.Count++
	r.ResponseRecorder.WriteHeader(statusCode)
}

func TestMustSendAsJSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		expResponse := `{"id":"abcdefghijkl","success":true}`
		expWriteHeaderCalls := 1

		payload := map[string]any{
			"success": true,
			"id":      "abcdefghijkl",
		}

		response.MustSendAsJSON(rr, payload, http.StatusOK, nil)

		require.Equal(expWriteHeaderCalls, rr.Count)
		assert.Equal(http.StatusOK, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResponse, rr.Body.String())
		}
	})

	t.Run("SuccessWithHeaders", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		expResp := `{"id":"abcdefghijkl","success":true}`
		expHeaders := http.Header{
			"Content-Type": []string{"application/json"},
		}
		expWriteHeaderCalls := 1

		payload := map[string]any{
			"success": true,
			"id":      "abcdefghijkl",
		}
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		response.MustSendAsJSON(rr, payload, http.StatusOK, headers)

		require.Equal(expWriteHeaderCalls, rr.Count)
		assert.Equal(http.StatusOK, rr.Code)
		assert.Equal(expHeaders, rr.Header())

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResp, rr.Body.String())
		}
	})

	// This should panic because the json encoder should
	// not be able to encode a Go channel
	t.Run("PanicForUnsupportedType", func(t *testing.T) {
		assert := assert.New(t)
		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		expWriteHeaderCalls := 1

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected a panic, but did not get one")
			}
			assert.Equal(http.StatusOK, rr.Code)
			assert.Equal(expWriteHeaderCalls, rr.Count, "got more than one WriteHeader calls, but expected one")
		}()

		invalidPayload := make(chan string)

		response.MustSendAsJSON(rr, invalidPayload, http.StatusOK, nil)
	})
}

func TestSendAsJSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		expResp := `{"id":"abcdefghijkl","success":true}`
		expWriteHeaderCalls := 1

		payload := map[string]any{
			"success": true,
			"id":      "abcdefghijkl",
		}

		response.SendAsJSON(rr, payload, http.StatusOK, nil)

		require.Equal(expWriteHeaderCalls, rr.Count)
		assert.Equal(http.StatusOK, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResp, rr.Body.String())
		}
	})

	t.Run("SuccessWithHeaders", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		expResp := `{"id":"abcdefghijkl","success":true}`
		expHeaders := http.Header{
			"Content-Type": []string{"application/json"},
		}
		expWriteHeaderCalls := 1

		payload := map[string]any{
			"success": true,
			"id":      "abcdefghijkl",
		}
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		response.SendAsJSON(rr, payload, http.StatusOK, headers)
		require.Equal(expWriteHeaderCalls, rr.Count)
		assert.Equal(http.StatusOK, rr.Code)
		assert.Equal(expHeaders, rr.Header())

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResp, rr.Body.String())
		}
	})

	t.Run("FailInternalServerError", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		expResp := `{"error":"json: unsupported type: chan string"}`
		expWriteHeaderCalls := 1

		invalidPayload := make(chan string)

		response.SendAsJSON(rr, invalidPayload, http.StatusOK, nil)
		require.Equal(expWriteHeaderCalls, rr.Count)

		assert.Equal(http.StatusInternalServerError, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResp, rr.Body.String())
		}
	})

	t.Run("SuccessNoSuperfluousResponseWriteHeaderCall", func(t *testing.T) {
		assert := assert.New(t)

		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		expWriteHeaderCalls := 1

		invalidPayload := make(chan string)

		response.SendAsJSON(rr, invalidPayload, http.StatusOK, nil)

		assert.Equal(http.StatusInternalServerError, rr.Code)
		assert.Equal(expWriteHeaderCalls, rr.Count, "got more than one WriteHeader calls, but expected one")
	})
}

func TestSendInternalServerErrorAsJSON(t *testing.T) {
	t.Run("SuccessErrorResponse", func(t *testing.T) {
		assert := assert.New(t)

		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		expWriteHeaderCalls := 1
		expCode := http.StatusInternalServerError
		expResp := `{"error":"an error has occurred"}`

		message := "an error has occurred"

		response.SendInternalServerErrorAsJSON(rr, message)

		assert.Equal(expWriteHeaderCalls, rr.Count, "got more than one WriteHeader calls, but expected one")
		assert.Equal(expCode, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResp, rr.Body.String())
		}
	})

	t.Run("SuccessNoSuperfluousResponseWriteHeaderCall", func(t *testing.T) {
		assert := assert.New(t)

		rr := &writeHeaderRecorder{ResponseRecorder: httptest.NewRecorder()}

		message := "an error has occurred"

		response.SendInternalServerErrorAsJSON(rr, message)

		expWriteHeaderCalls := 1
		assert.Equal(expWriteHeaderCalls, rr.Count, "got more than one WriteHeader calls, but expected one")
	})
}
