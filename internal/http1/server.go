package http1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
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
	reader := bufio.NewReader(c)
	ParseRequest(reader)

	//fmt.Printf("%s", content)

	// if err != nil {
	// 	fmt.Printf("Failed to parse")
	// 	log.Fatal(err)
	// }

	/*
			Functions that might be useful:
			* ParseRequest()
		  * mux.ServeHTTP()
	*/
	return nil
}
