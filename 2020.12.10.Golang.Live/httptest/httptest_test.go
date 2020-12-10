package httptest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// start_Test_handler OMIT
func Test_handler(t *testing.T) {
	wantBody := `{"path":"/aTest"}`

	w := httptest.NewRecorder()                                                // HLrecorder
	r := httptest.NewRequest(http.MethodGet, "https://example.com/aTest", nil) // HLrecorder

	handler(w, r)

	response := w.Result() // HLresponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("could not read result body: %v", err)
	}

	if response.StatusCode != http.StatusOK { // HLresponse
		t.Errorf("StatusCode: got %d, want 200", response.StatusCode)
	}
	if wantBody != string(body) {
		t.Errorf("Body: got %s, want %s", string(body), wantBody)
	}
}

// end_Test_handler OMIT

// start_TestNewFetcher OMIT
func TestNewFetcher(t *testing.T) {
	want := "42"
	// start_testServer OMIT
	ts := httptest.NewServer( // HLserver
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // HLserver
			id := strings.Split(r.URL.Path, "/")[1]
			_, err := strconv.Atoi(id)

			if err != nil {
				t.Errorf("id is not a number: %v", err) // HLassertions
			}
			if id != want {
				t.Errorf("want: %s, got: %s", want, id) // HLassertions
			}
		}))
	defer ts.Close() // HLserver
	// end_testServer OMIT

	// start_NewFetcher OMIT
	f := NewFetcher(ts.URL) // HLNewFetcher
	// end_NewFetcher OMIT
	_, err := f(want)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// end_TestNewFetcher OMIT
