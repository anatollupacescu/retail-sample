package acceptance

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/httpclient"
)

var (
	apiURL  = flag.String("apiURL", "", "api server URL")
	timeout = 100 * time.Millisecond //TODO pass as flag
)

func Post() func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")

		return httpClient.Post(*apiURL, body, headers)
	}
}

func Patch(id int) func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")

		resourceURL := fmt.Sprintf("%s/%d", *apiURL, id)

		return httpClient.Patch(resourceURL, body, headers)
	}
}

func Get(id... int) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Accept", "application/json")

		var resourceURL = *apiURL

		if len(id) == 1 {
			resourceURL = fmt.Sprintf("%s/%d", *apiURL, id[0])
		}

		return httpClient.Get(resourceURL, headers)
	}
}
