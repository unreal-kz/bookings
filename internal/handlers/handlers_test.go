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

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{name: "home", url: "/", method: "GET", params: []postData{}, expectedStatusCode: http.StatusOK},
	{name: "about", url: "/about", method: "GET", params: []postData{}, expectedStatusCode: http.StatusOK},
	{name: "sa", url: "/search-availability", method: "GET", params: []postData{}, expectedStatusCode: http.StatusOK},
	{name: "ms", url: "/majors-suite", method: "GET", params: []postData{}, expectedStatusCode: http.StatusOK},
	{name: "gq", url: "/generals-quarters", method: "GET", params: []postData{}, expectedStatusCode: http.StatusOK},
	{name: "mr", url: "/make-reservation", method: "GET", params: []postData{}, expectedStatusCode: http.StatusOK},
	{name: "rs", url: "/reservation-summary", method: "GET", params: []postData{}, expectedStatusCode: http.StatusOK},
	{name: "contact", url: "/contact", method: "GET", params: []postData{}, expectedStatusCode: http.StatusOK},
	{name: "post_sa", url: "/search-availability", method: "POST", params: []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-04"},
	}, expectedStatusCode: http.StatusOK},
	{name: "post_sa_json", url: "/search-availability-json", method: "POST", params: []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-04"},
	}, expectedStatusCode: http.StatusOK},
	{name: "post_mr", url: "/make-reservation", method: "POST", params: []postData{
		{key: "first_name", value: "Bob"},
		{key: "last_name", value: "Uki"},
		{key: "email", value: "uki@gmail.com"},
		{key: "phone", value: "555-555"},
	}, expectedStatusCode: http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			var values = url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
