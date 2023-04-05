package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ricxi/flat-list/shared/response"
	"github.com/stretchr/testify/assert"
)

func TestHandlePanic(t *testing.T) {
	t.Run("simplePanic", func(t *testing.T) {
		assert := assert.New(t)

		panicker := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("this handler panics")
		})
		handlePanic := response.HandlePanic(panicker)

		expCode := http.StatusInternalServerError
		expResp := `{"error": "an internal server error has occurred"}`

		r := httptest.NewRequest("", "/", nil)
		rr := httptest.NewRecorder()

		handlePanic.ServeHTTP(rr, r)

		assert.Equal(expCode, rr.Code)
		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResp, rr.Body.String())
		}
	})

	t.Run("handlePanicFromMustSendJSON", func(t *testing.T) {
		assert := assert.New(t)

		panicker := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response.MustSendJSON(w, make(chan string), http.StatusInternalServerError, nil)
		})
		handlePanic := response.HandlePanic(panicker)

		expCode := http.StatusInternalServerError
		expResp := `{"error": "an internal server error has occurred"}`

		r := httptest.NewRequest("", "/", nil)
		rr := httptest.NewRecorder()

		handlePanic.ServeHTTP(rr, r)

		assert.Equal(expCode, rr.Code)
		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResp, rr.Body.String())
		}
	})
}
