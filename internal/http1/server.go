package http1

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func ListenAndServe(addr string, handler http.Handler) error {

	ln, err := net.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer ln.Close()

	if handler == nil {
		handler = http.DefaultServeMux
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go HandleConnection(conn, handler)
	}
}

func HandleConnection(c net.Conn, mux http.Handler) error {

	defer c.Close()

	reader := bufio.NewReader(c)

	for {
		// TODO : What if we want to consume the client request body
		req, err := ParseRequest(reader)
		if err != nil {
			// client closes connection
			if errors.Is(err, io.EOF) {
				break
			}

			var httpErr *HTTPError
			code := http.StatusBadRequest
			if errors.As(err, &httpErr) {
				code = httpErr.StatusCode
			}

			w := NewResponse(c)
			w.Header().Set("Connection", "close")
			http.Error(w, http.StatusText(code), code)
			w.Finish()

			break
		}
		if req == nil {
			break
		}
		fmt.Printf("[+] Request received : %s\n", req.Method)

		// Print the request (just for debug)
		// dump, _ := httputil.DumpRequest(req, true)
		// fmt.Println(string(dump))

		w := NewResponse(c)
		mux.ServeHTTP(w, req)

		w.Finish()

		if req.Body != nil {
			io.Copy(io.Discard, req.Body)
			req.Body.Close()
		}

		keepAlive := true
		connHeader := strings.ToLower(req.Header.Get("Connection"))
		if connHeader == "keep-alive" {
			keepAlive = true
		} else if connHeader == "close" || req.ProtoMinor == 0 {
			keepAlive = false
		}

		if !keepAlive {
			break
		}
	}

	return nil
}
