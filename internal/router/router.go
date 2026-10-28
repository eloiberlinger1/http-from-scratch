package router

import (
	"httpServer/internal/handlers"
	"httpServer/internal/mux1"
)

func SetupRouter() *mux1.ServeMux {
	mux := mux1.NewServeMux()

	handlers.RegisterRootHandler(mux)
	handlers.RegisterHelloHandler(mux)

	return mux
}
