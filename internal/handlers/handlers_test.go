package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"genq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majs", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"makeres", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"ressummary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "POST", []postData{
		{"start", "2023-10-01"},
		{"end", "2023-10-02"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
		{"start", "2023-10-01"},
		{"end", "2023-10-02"},
	}, http.StatusOK},
	{"post-make-res", "/make-reservation", "POST", []postData{
		{"first_name", "John"},
		{"last_name", "Doe"},
		{"email", "me@email.org"},
		{"phone", "1234567890"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes) // ts standas for test server
	defer ts.Close()

	for _, e := range theTest {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Logf("Error making GET request: %v", err)
				t.Fatal()
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For test %s, Expected status code %d, got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			// POST request
			form := url.Values{}
			for _, p := range e.params {
				form.Add(p.key, p.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, form)
			if err != nil {
				t.Logf("Error making POST request: %v", err)
				t.Fatal()
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For test %s, Expected status code %d, got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
