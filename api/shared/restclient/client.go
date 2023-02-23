package restclient

import "net/http"

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func Get(url string) (resp *http.Response, err error) {
	response, err := Client.Get(url)
	if err != nil {
		return nil, err
	}

	return response, err
}
