package handlers

import (
	"httpServer/internal/mux1"
	"net/http"
	"strconv"
)

func RegisterRootHandler(mux *mux1.ServeMux) {

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := "root endpoint"

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

}

func RegisterHelloHandler(mux *mux1.ServeMux) {

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		data := "HELLO "

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

}
