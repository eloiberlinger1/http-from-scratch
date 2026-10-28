package main

import (
	"httpServer/internal/http1"
	"httpServer/internal/router"
)

const addr = ":8080"

func main() {

	mux := router.SetupRouter()

	if err := http1.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}

}
