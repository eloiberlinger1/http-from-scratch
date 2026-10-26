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
	reader := bufio.NewReader(c)
	req, err := ParseRequest(reader)
	if err != nil {
		return err
	}

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Printf("Erreur lors du dump : %v\n", err)
		return err
	}

	fmt.Println("--- REQUÊTE RECONSTRUITE ---")
	fmt.Println(string(dump))

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
