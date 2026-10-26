package http1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
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
			log.Fatal(err)
		}
		go HandleConnection(conn, handler)
	}
}

func HandleConnection(c net.Conn, mux http.Handler) error {

	defer c.Close()

	reader := bufio.NewReader(c)
	req, err := ParseRequest(reader)

	if err != nil {
		fmt.Printf("Failed to parse")
		log.Fatal(err)
	}

	// Print the request (just for debug)
	dump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(dump))

	// response
	w := NewResponse(c)

	mux.ServeHTTP(w, req)

	return nil
}
