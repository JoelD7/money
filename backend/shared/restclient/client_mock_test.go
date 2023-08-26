package restclient

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	c := require.New(t)

	client := NewMockRestClient()

	t.Run("Single url", func(t *testing.T) {
		urlOne := "https://localhost"

		err := client.AddMockedResponseFromFile("samples/response.json", urlOne, MethodGET)
		c.Nil(err)

		response, err := client.Get(urlOne)
		c.Nil(err)

		body, err := readResponseBody(response)
		c.Nil(err)
		c.Equal(`{"key": "value"}`, body)
	})

	t.Run("Multiple urls", func(t *testing.T) {
		urlOne := "https://localhost"
		urlTwo := "http://example.com"

		err := client.AddMockedResponseFromFile("samples/response.json", urlOne, MethodGET)
		c.Nil(err)

		err = client.AddMockedResponseFromFile("samples/response_two.json", urlTwo, MethodGET)
		c.Nil(err)

		response, err := client.Get(urlOne)
		c.Nil(err)

		body, err := readResponseBody(response)
		c.Nil(err)
		c.Equal(`{"key": "value"}`, body)

		response, err = client.Get(urlTwo)
		c.Nil(err)

		body, err = readResponseBody(response)
		c.Nil(err)
		c.Equal(`{"key2": "value2"}`, body)
	})

	t.Run("Mocked response by method overrides mocked by method and url", func(t *testing.T) {
		urlOne := "https://localhost"
		urlTwo := "http://example.com"

		err := client.AddMockedResponseFromFile("samples/response.json", urlOne, MethodGET)
		c.Nil(err)

		response, err := client.Get(urlOne)
		c.Nil(err)

		body, err := readResponseBody(response)
		c.Nil(err)
		c.Equal(`{"key": "value"}`, body)

		err = client.AddMockedResponseFromFile("samples/response_two.json", urlTwo, MethodGET)
		c.Nil(err)

		response, err = client.Get(urlTwo)
		c.Nil(err)

		body, err = readResponseBody(response)
		c.Nil(err)
		c.Equal(`{"key2": "value2"}`, body)

		err = client.AddMockedResponseFromFileNoUrl("samples/response.json", MethodGET)
		c.Nil(err)

		response, err = client.Get(urlOne)
		c.Nil(err)

		body, err = readResponseBody(response)
		c.Nil(err)
		c.Equal(`{"key": "value"}`, body)

		response, err = client.Get(urlTwo)
		c.Nil(err)

		body, err = readResponseBody(response)
		c.Nil(err)
		c.Equal(`{"key": "value"}`, body)
	})
}

func readResponseBody(response *http.Response) (strBody string, err error) {
	defer func() {
		//Here I'm purposely ignoring the error from the Close() function if other part of the code already returned
		//an error, as I consider the latter more relevant to the application than the former.
		tempErr := response.Body.Close()
		if tempErr != nil && err != nil {
			return
		}

		if tempErr != nil {
			err = tempErr
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	response.Body = io.NopCloser(bytes.NewReader(body))

	return string(body), nil
}
