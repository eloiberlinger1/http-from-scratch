package main

import (
	"fmt"
	"httpServer/internal/http1"
	"net/http"
	"os"
	"strconv"
)

const addr = ":8080"

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		fmt.Fprintf(w, "hello")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./internal/src/index.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 not found"))
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	if err := http1.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}

}
