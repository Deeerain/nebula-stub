package main

import (
	"fmt"
	"net/http"
)

var (
	address = "localhost:8080"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	})

	http.ListenAndServe(address, nil)
}
