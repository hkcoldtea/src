package main

import (
	"log"
	"net/http"
	"os"

	"controller"
)

func main() {
	if err := controller.Run(); err != nil && err != http.ErrServerClosed {
		// unexpected error. port in use?
		log.Fatalf("ListenAndServe() Err: %v", err)
		os.Exit(1)
	}
}
