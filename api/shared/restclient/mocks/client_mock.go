package mocks

import "net/http"

type MockClient struct{}

var (
	GetFunction func(url string) (*http.Response, error)
)

func (m *MockClient) Get(url string) (*http.Response, error) {
	return GetFunction(url)
}
