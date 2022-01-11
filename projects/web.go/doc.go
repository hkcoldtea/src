/*
 Package Document

# This script tests that running go mod with
# GO111MODULE=off

 I ran this command export GO111MODULE="off" and that worked for me.

 Sample code:

package main

import (
	"github.com/hkcoldtea/src/projects/web.go/server"
)

func main() {
	s := server.InitServer()

	s.Get("/stop", func(c *server.Context) {
		c.Shutdown()
	})

	s.Run("localhost:9999")
}
*/
package main
