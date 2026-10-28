package main

import (
	"httpServer/internal/http1"
	"net/http"
	"os"
	"strconv"
)

const addr = ":8080"

func main() {

	mux := http.NewServeMux()

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
