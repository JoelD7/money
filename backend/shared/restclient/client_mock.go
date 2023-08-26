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
	mockedResponses         map[Method]map[string]*http.Response
	mockedResponsesByMethod map[Method]*http.Response
}

func NewMockRestClient() *MockClient {
	return &MockClient{
		mockedResponses:         make(map[Method]map[string]*http.Response),
		mockedResponsesByMethod: make(map[Method]*http.Response),
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

	m.mockedResponsesByMethod[method] = nil

	return nil
}

// AddMockedResponseFromFileNoUrl uses the contents of the file at <path> to mock the response for the specified method.
// Keep in mind that this function will mock the responses of any request of type <method>, disregarding the url
func (m *MockClient) AddMockedResponseFromFileNoUrl(path string, method Method) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	r := io.NopCloser(bytes.NewReader(data))

	m.mockedResponsesByMethod[method] = &http.Response{
		StatusCode:    http.StatusOK,
		Body:          r,
		ContentLength: 1,
	}

	m.mockedResponses[method] = nil

	return nil
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	if m.mockedResponsesByMethod[MethodGET] != nil {
		return m.mockedResponsesByMethod[MethodGET], nil
	}

	if m.mockedResponses[MethodGET] == nil && m.mockedResponses[MethodGET][url] == nil {
		r := io.NopCloser(bytes.NewReader([]byte{}))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       r,
		}, nil
	}

	return m.mockedResponses[MethodGET][url], nil
}
