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

func RegisterEchoHandler(mux *mux1.ServeMux) {

	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		data := "echo endpoint"

		// just test to consume body for the test TestHandleConnection_CornerCases/body_shorter_than_Content-Length
		// if r.ContentLength > 0 {
		// 	// io.ReadAll alloue par morceaux, ou io.Copy(io.Discard, ...) qui n'alloue rien
		// 	_, err := io.ReadAll(r.Body)
		// 	if err != nil {
		// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		// 		return
		// 	}
		// }

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

}
