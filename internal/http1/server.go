package http1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type HTTPError struct {
	StatusCode int
}

func (e HTTPError) WriteError(c net.Conn) {
	statusText := http.StatusText(e.StatusCode)
	body := fmt.Sprintf("%d %s\n", e.StatusCode, statusText)

	resp := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s",
		e.StatusCode,
		statusText,
		len(body),
		body,
	)
	c.Write([]byte(resp))
}

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
		// return the corresponding error
		fmt.Printf("[+] Request received : %s\n", http.StatusText(err.StatusCode))
		err.WriteError(c)
		return nil
	}
	fmt.Printf("[+] Request received : %s\n", req.Method)

	// Print the request (just for debug)
	// dump, _ := httputil.DumpRequest(req, true)
	// fmt.Println(string(dump))

	// prepare response object
	w := NewResponse(c)
	mux.ServeHTTP(w, req)

	return nil
}
