package http1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
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

	if err := c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}

	reader := bufio.NewReader(c)
	req, err := ParseRequest(reader)
	if err != nil {
		fmt.Printf("Failed to parse: %v", err)
		log.Fatal(err)
	}

	// Print the request (just for debug)
	// dump, _ := httputil.DumpRequest(req, true)
	// fmt.Println(string(dump))
	fmt.Printf("1 Request received : %s\n", req.Method)

	// prepare response object
	w := NewResponse(c)

	mux.ServeHTTP(w, req)

	return nil
}
