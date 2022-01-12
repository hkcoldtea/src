package util

import (
	"os"
	"testing"
)

func TestPamAuthenticate(t *testing.T) {
	username := os.Getenv("USER")
	password := os.Getenv("PASS")
	t.Log(username)
	t.Log(password)
	ret := PamAuthenticate(username, password)
	t.Log(ret)
	if ret != 1 {
		t.Fatal(ret)
	}
}
