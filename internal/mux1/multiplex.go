package mux1

import (
	"net/http"
)

type Route struct {
	Pattern string
	Handler func(w http.ResponseWriter, r *http.Request)
}

type ServeMux struct {
	Routes []Route
}

func (m *ServeMux) HandleFunc(ep string, f func(w http.ResponseWriter, r *http.Request)) {
	m.Routes = append(m.Routes, Route{
		Pattern: ep,
		Handler: f,
	})
}

func NewServeMux() *ServeMux {
	return &ServeMux{}
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range m.Routes {
		if route.Pattern == r.URL.Path {
			route.Handler(w, r)
			return
		}
	}
	http.NotFound(w, r)

}
