package controllers

import (
	"fmt"
	"net/http"
)

// HelloWorld to test mux setup
var HelloWorld = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}
