package main

import (
	"fmt"
	"httpServer/internal/http1"
	"net/http"
)

const addr = ":8080"

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		fmt.Fprintf(w, "hello")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	if err := http1.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}

}
