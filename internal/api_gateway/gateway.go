// internal/api_gateway/gateway.go

package api_gateway

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter initializes and returns a new Gorilla mux router.
func NewRouter() *mux.Router {
	return mux.NewRouter()
}

// RunServer runs the HTTP server with the specified handler and address.
func RunServer(handler http.Handler, addr string) error {
	return http.ListenAndServe(addr, handler)
}
