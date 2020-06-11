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

func Post(resourceName string, id ...int) func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")

		resourceURL := fmt.Sprintf("%s/%s", *apiURL, resourceName)

		if len(id) > 0 {
			resourceURL = fmt.Sprintf("%s/%d", resourceURL, id[0])
		}

		return httpClient.Post(resourceURL, body, headers)
	}
}

func Patch(resourceName string, id int) func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")

		resourceURL := fmt.Sprintf("%s/%s/%d", *apiURL, resourceName, id)

		return httpClient.Patch(resourceURL, body, headers)
	}
}

func Get(resourceName string, id ...int) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
		)
		headers := http.Header{}
		headers.Set("Accept", "application/json")

		var resourceURL = fmt.Sprintf("%s/%s", *apiURL, resourceName)

		if len(id) == 1 {
			resourceURL = fmt.Sprintf("%s/%d", resourceURL, id[0])
		}

		return httpClient.Get(resourceURL, headers)
	}
}
