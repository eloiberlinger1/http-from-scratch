package main

import (
	"httpServer/internal/http1"
	"httpServer/internal/http1/middlewares"
	"httpServer/internal/router"
)

const addr = ":8080"

func main() {

	mux := router.SetupRouter()

	handlerWithMetrics := middlewares.MetricsMiddleware(mux)

	if err := http1.ListenAndServe(addr, handlerWithMetrics); err != nil {
		panic(err)
	}

}
