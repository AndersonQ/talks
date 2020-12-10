package httptest

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// start_handler OMIT
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"path":"%s"}`, r.URL.Path)
}

// end_handler OMIT

// start_NewFetcher OMIT
func NewFetcher(baseURL string) func(id string) (string, error) {
	return func(id string) (string, error) {
		// start_fetcher OMIT
		resp, err := http.Get(baseURL + "/" + id)
		if err != nil {
			return "", fmt.Errorf("could not fetch something: %w", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("could not read response body: %w", err)
		}

		return string(body), nil
		// end_fetcher OMIT
	}
}

// end_NewFetcher OMIT
