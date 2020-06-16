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

func PostToID(resourceName string, id int) func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient, headers := httpClient()
		resourceURL := fmt.Sprintf("%s/%s/%d", *apiURL, resourceName, id)
		return httpClient.Post(resourceURL, body, headers)
	}
}

func Post(resourceName string) func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient, headers := httpClient()
		resourceURL := fmt.Sprintf("%s/%s", *apiURL, resourceName)
		return httpClient.Post(resourceURL, body, headers)
	}
}

func Patch(resourceName string, id int) func(io.Reader) (*http.Response, error) {
	return func(body io.Reader) (*http.Response, error) {
		httpClient, headers := httpClient()
		entityURL := fmt.Sprintf("%s/%s/%d", *apiURL, resourceName, id)
		return httpClient.Patch(entityURL, body, headers)
	}
}

func List(resourceName string) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		httpClient, headers := httpClient()
		var resourceURL = fmt.Sprintf("%s/%s", *apiURL, resourceName)
		return httpClient.Get(resourceURL, headers)
	}
}

func Get(resourceName string, id int) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		httpClient, headers := httpClient()
		var entityURL = fmt.Sprintf("%s/%s/%d", *apiURL, resourceName, id)
		return httpClient.Get(entityURL, headers)
	}
}

func httpClient() (*httpclient.Client, http.Header) {
	c := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
	)
	headers := http.Header{}
	headers.Set("Accept", "application/json")

	return c, headers
}
