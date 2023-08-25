package restclient

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

var (
	MethodGET Method = http.MethodGet
)

type Method string

type MockClient struct {
	// List of responses by url by method. This map reads: responses of X url for Y method.
	mockedResponses map[Method]map[string]*http.Response
}

func NewMockRestClient() *MockClient {
	return &MockClient{
		mockedResponses: make(map[Method]map[string]*http.Response),
	}
}

// AddMockedResponseFromFile uses the contents of the file at <path> to mock the response for the specified method and url.
func (m *MockClient) AddMockedResponseFromFile(path string, url string, method Method) error {
	if m.mockedResponses[method] == nil {
		m.mockedResponses[method] = map[string]*http.Response{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	r := io.NopCloser(bytes.NewReader(data))

	m.mockedResponses[method][url] = &http.Response{
		StatusCode: http.StatusOK,
		Body:       r,
	}

	return nil
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	if m.mockedResponses[MethodGET] == nil || m.mockedResponses[MethodGET][url] == nil {
		r := io.NopCloser(bytes.NewReader([]byte{}))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       r,
		}, nil
	}

	return m.mockedResponses[MethodGET][url], nil
}
