package apigateway

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewJSONResponse(t *testing.T) {
	resetOriginsMap := func() {
		allowedOriginsMap = map[string]struct{}{}
	}

	t.Run("Success with struct body", func(t *testing.T) {
		resetOriginsMap()
		req := &Request{Headers: map[string]string{}}
		body := map[string]string{"message": "success"}

		resp := req.NewJSONResponse(http.StatusOK, body)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Headers["Content-Type"])

		var respBody map[string]string
		err := json.Unmarshal([]byte(resp.Body), &respBody)
		assert.NoError(t, err)
		assert.Equal(t, "success", respBody["message"])
	})

	t.Run("Success with string body", func(t *testing.T) {
		resetOriginsMap()
		req := &Request{Headers: map[string]string{}}
		body := "raw string content"

		resp := req.NewJSONResponse(http.StatusCreated, body)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "raw string content", resp.Body)
	})

	t.Run("CORS allowed origin", func(t *testing.T) {
		resetOriginsMap()
		t.Setenv("CORS_ORIGIN", "https://example.com;https://trusted.com")

		req := &Request{
			Headers: map[string]string{
				"origin": "https://example.com",
			},
		}

		resp := req.NewJSONResponse(http.StatusOK, "body")

		assert.Equal(t, "https://example.com", resp.Headers["Access-Control-Allow-Origin"])
		assert.Equal(t, "true", resp.Headers["Access-Control-Allow-Credentials"])
	})

	t.Run("CORS disallowed origin", func(t *testing.T) {
		resetOriginsMap()
		t.Setenv("CORS_ORIGIN", "https://trusted.com")

		req := &Request{
			Headers: map[string]string{
				"origin": "https://malicious.com",
			},
		}

		resp := req.NewJSONResponse(http.StatusOK, "body")

		_, ok := resp.Headers["Access-Control-Allow-Origin"]
		assert.False(t, ok, "Should not return Access-Control-Allow-Origin for untrusted origin")
	})

	t.Run("CORS wildcard origin", func(t *testing.T) {
		resetOriginsMap()
		t.Setenv("CORS_ORIGIN", "*")

		req := &Request{
			Headers: map[string]string{
				"origin": "https://anywhere.com",
			},
		}

		resp := req.NewJSONResponse(http.StatusOK, "body")

		_, ok := resp.Headers["Access-Control-Allow-Origin"]
		assert.True(t, ok)
	})

	t.Run("Custom headers", func(t *testing.T) {
		resetOriginsMap()
		req := &Request{Headers: map[string]string{}}
		customHeader := Header{Key: "X-Trace-ID", Value: "12345"}

		resp := req.NewJSONResponse(http.StatusOK, "body", customHeader)

		assert.Equal(t, "12345", resp.Headers["X-Trace-ID"])
	})

	t.Run("Marshal error handling", func(t *testing.T) {
		resetOriginsMap()
		req := &Request{Headers: map[string]string{}}

		// Create a body that fails JSON marshaling (e.g., a channel)
		body := map[string]interface{}{
			"data": make(chan int),
		}

		// When marshal fails, NewJSONResponse calls NewErrorResponse
		// NewErrorResponse will likely fall through to ErrInternalError (500)
		resp := req.NewJSONResponse(http.StatusOK, body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var respErr Error
		err := json.Unmarshal([]byte(resp.Body), &respErr)
		assert.NoError(t, err)
		assert.NotEmpty(t, respErr.Message)
	})

	t.Run("Nil body is empty response", func(t *testing.T) {
		resetOriginsMap()
		req := &Request{Headers: map[string]string{}}

		resp := req.NewJSONResponse(http.StatusOK, nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Headers["Content-Type"])

		assert.Empty(t, resp.Body)
	})
}
