package drudapi

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func getTestServer(code int, body string) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	return server
}

func TestGetClient200(t *testing.T) {
	expectedResp := `{
    "_updated": "Mon, 23 May 2016 20:23:34 GMT",
    "name": "1fee",
    "email": "me@there.com",
    "phone": "123-123-1234",
    "_links": {
        "self": {
            "href": "client/1fee",
            "title": "Client"
        },
        "parent": {
            "href": "/",
            "title": "home"
        },
        "collection": {
            "href": "client",
            "title": "client"
        }
    },
    "_created": "Mon, 23 May 2016 20:23:34 GMT",
    "_id": "574366c6e2638a001f430115",
    "_etag": "9906d3a8584f0fabbd96451013b38d20fff5f5d3"
}`
	server := getTestServer(200, expectedResp)
	defer server.Close()

	r := Request{
		Host: server.URL,
		Auth: &Credentials{
			AdminToken: os.Getenv("GITHUB_TOKEN"),
		},
	}

	c := &Client{
		Name: "1fee",
	}

	err := r.Get(c)
	if err != nil {
		log.Fatal(err)
	}

	expect(t, c.Name, "1fee")
	expect(t, c.Email, "me@there.com")
	expect(t, c.Phone, "123-123-1234")
	expect(t, c.Created, "Mon, 23 May 2016 20:23:34 GMT")
	expect(t, c.ID, "574366c6e2638a001f430115")
	expect(t, c.Etag, "9906d3a8584f0fabbd96451013b38d20fff5f5d3")

}

func TestPostClient200(t *testing.T) {
	expectedResp := `{
    "_updated": "Mon, 23 May 2016 20:23:34 GMT",
    "name": "whosit",
    "email": "what@where.com",
    "phone": "123-123-1234"
}`
	server := getTestServer(200, expectedResp)
	defer server.Close()

	r := Request{
		Host: server.URL,
		Auth: &Credentials{
			AdminToken: os.Getenv("GITHUB_TOKEN"),
		},
	}

	c := &Client{
		Name:  "whosit",
		Email: "what@where.com",
	}

	err := r.Post(c)
	if err != nil {
		log.Fatal(err)
	}

	expect(t, c.Name, "whosit")
	expect(t, c.Email, "what@where.com")

}
