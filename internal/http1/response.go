package http1

import (
	"fmt"
	"net"
	"net/http"
)

type Response struct {
	Connection net.Conn
	header     http.Header
	StatusCode int
	ProtoMajor int
	ProtoMinor int
	HeaderSent bool
}

func NewResponse(conn net.Conn) *Response {

	return &Response{
		Connection: conn,
		header:     make(http.Header),
		StatusCode: 200,
		HeaderSent: false,
	}
}

func (r *Response) Header() http.Header {

	return r.header
}

func (r *Response) Write(data []byte) (n int, err error) {

	if !r.HeaderSent {
		r.WriteHeader(200)
	}

	return r.Connection.Write(data)

}

func (r *Response) WriteHeader(statusCode int) {
	// implement statuscode
	//
	if r.HeaderSent {
		return
	}
	r.StatusCode = statusCode

	statusText := http.StatusText(statusCode)
	fmt.Fprintf(r.Connection, "HTTP/1.1 %d %s\r\n", statusCode, statusText)

	for k, values := range r.header {
		for _, v := range values {
			fmt.Fprintf(r.Connection, "%s: %s\r\n", k, v)
		}
	}

	fmt.Fprintf(r.Connection, "\r\n")
	r.HeaderSent = true

}
