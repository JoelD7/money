package restclient

import "net/http"

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type RestClient struct {
	client *http.Client
}

func New() *RestClient {
	return &RestClient{
		client: new(http.Client),
	}
}

func (r *RestClient) Get(url string) (resp *http.Response, err error) {
	response, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}

	return response, err
}
