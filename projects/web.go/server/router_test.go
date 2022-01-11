package server

import (
	"testing"
)

func newTestRouter() *Router {
	r := InitRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	url, pathParams := r.getRoute("GET", "/hello/golang")

	if url == "" {
		t.Fatal("url shouldn't be empty")
	}

	if url != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if pathParams["name"] != "golang" {
		t.Fatal("name should be equal to 'golang'")
	}

	t.Logf("matched path: %s, params['name']: %s\n", url, pathParams["name"])
}

func TestGetRouter2(t *testing.T) {
	r := newTestRouter()
	n1, ps1 := r.getRoute("GET", "/assets/file1.txt")
	ok1 := n1 == "/assets/*filepath" && ps1["filepath"] == "file1.txt"
	if !ok1 {
		t.Fatal("n1 shoule be /assets/*filepath & filepath shoule be file1.txt")
	}

	n2, ps2 := r.getRoute("GET", "/assets/css/test.css")
	ok2 := n2 == "/assets/*filepath" && ps2["filepath"] == "css/test.css"
	if !ok2 {
		t.Fatal("n2 shoule be /assets/*filepath & filepath shoule be css/test.css")
	}
}
