package request_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ricxi/flat-list/shared/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// output is used to mock a
// destination for ParseJSON
// to decode data into
type output struct {
	Message string `json:"message"`
}

func TestParseJSON(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		assert := assert.New(t)

		expectedOutput := output{
			Message: "an apple a day",
		}

		input := `{"message":"an apple a day"}`

		r := httptest.NewRequest("", "/", strings.NewReader(input))
		r.Header.Set("Content-Type", "application/json")

		var actualOutput output
		err := request.ParseJSON(r, &actualOutput)

		require.NoError(t, err)

		if assert.NotEmpty(actualOutput) {
			assert.Equal(expectedOutput.Message, actualOutput.Message)
		}
	})

	t.Run("InvalidJSONSyntaxError", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		input := `{"message":"an apple a day",}`

		r := httptest.NewRequest("", "/", strings.NewReader(input))
		r.Header.Set("Content-Type", "application/json")

		var actualOutput output
		err := request.ParseJSON(r, &actualOutput)

		require.Error(err)
		require.Empty(actualOutput)
		assert.IsType(&json.SyntaxError{}, err)
	})

	t.Run("InvalidJSONUnexpectedEOF", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		input := `{"message":"an apple a day"`

		r := httptest.NewRequest("", "/", strings.NewReader(input))
		r.Header.Set("Content-Type", "application/json")

		var actualOutput output
		err := request.ParseJSON(r, &actualOutput)

		require.Error(err)
		require.Empty(actualOutput)
		assert.Equal(io.ErrUnexpectedEOF, err)
	})
}

func BenchmarkParseJSON(b *testing.B) {
	input := `{"message":"an apple a day"}`

	r := httptest.NewRequest("", "/", strings.NewReader(input))

	r.Header.Set("Content-Type", "application/json")

	for n := 0; n < b.N; n++ {
		var actualOutput output
		_ = request.ParseJSON(r, &actualOutput)
	}
}
