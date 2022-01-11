package server

import "testing"

func dummy(*Context) {
}

func TestNestedGroup(t *testing.T) {
	server := InitServer()
	group := server.SetGroup("/pre")
	{
		group.Get("/:path/bbb/:ccc", dummy)
	}
}
