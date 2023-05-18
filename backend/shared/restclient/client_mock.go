package restclient

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

var (
	MethodGET Method = http.MethodGet

	// List of responses by url by method. This map reads: responses of X url for Y method.
	mockedResponses = map[Method]map[string]*http.Response{}
)

type Method string

type MockClient struct{}

func InitMockClient() {
	Client = &MockClient{}
}

// AddMockedResponseFromFile uses the contents of the file at <path> to mock the response for the specified method and url.
func AddMockedResponseFromFile(path string, url string, method Method) error {
	if mockedResponses[method] == nil {
		mockedResponses[method] = map[string]*http.Response{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	r := io.NopCloser(bytes.NewReader(data))

	mockedResponses[method][url] = &http.Response{
		StatusCode: http.StatusOK,
		Body:       r,
	}

	return nil
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	if mockedResponses[MethodGET] == nil || mockedResponses[MethodGET][url] == nil {
		r := io.NopCloser(bytes.NewReader([]byte{}))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       r,
		}, nil
	}

	return mockedResponses[MethodGET][url], nil
}
