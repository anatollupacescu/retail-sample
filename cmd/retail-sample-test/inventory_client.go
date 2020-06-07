package inv

import (
	"io"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/httpclient"
)

func Post(url string, timeout time.Duration) func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")

		return httpClient.Post(url, body, headers)
	}
}

func Patch(url string, timeout time.Duration) func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")

		return httpClient.Patch(url, body, headers)
	}
}

func Get(url string, timeout time.Duration) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Accept", "application/json")

		return httpClient.Get(url, headers)
	}
}
