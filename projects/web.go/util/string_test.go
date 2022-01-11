package util

import (
	"reflect"
	"testing"
)

func TestParseUrl(t *testing.T) {
	ok := reflect.DeepEqual(ParseUrl("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(ParseUrl("/p/*"), []string{"p", "*"})
	t.Log(ParseUrl("/p/*name/*"))
	ok = ok && reflect.DeepEqual(ParseUrl("/p/*name/*"), []string{"p", "*name", "*"})
	if !ok {
		t.Fatal("test ParseUrl failed")
	}
}
