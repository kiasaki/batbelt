package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type GET interface {
	Path() string
	GET(http.ResponseWriter, *http.Request)
}

type POST interface {
	Path() string
	POST(http.ResponseWriter, *http.Request)
}

type PUT interface {
	Path() string
	PUT(http.ResponseWriter, *http.Request)
}

type DELETE interface {
	Path() string
	DELETE(http.ResponseWriter, *http.Request)
}

type HEAD interface {
	Path() string
	HEAD(http.ResponseWriter, *http.Request)
}

func handlerMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	SetMethodNotAllowedResponse(w)
}

// Registers handlers for all supported methods of an endopoint in
// a gorilla/mux router
func RegisterEnpointToRouter(router *mux.Router, endpoint interface{}) {
	if e, ok := endpoint.(GET); ok {
		router.HandleFunc(e.Path(), e.GET).Methods("GET")
	}
	if e, ok := endpoint.(POST); ok {
		router.HandleFunc(e.Path(), e.POST).Methods("POST")
	}
	if e, ok := endpoint.(PUT); ok {
		router.HandleFunc(e.Path(), e.PUT).Methods("PUT")
	}
	if e, ok := endpoint.(DELETE); ok {
		router.HandleFunc(e.Path(), e.DELETE).Methods("DELETE")
	}
	if e, ok := endpoint.(HEAD); ok {
		router.HandleFunc(e.Path(), e.HEAD).Methods("HEAD")
	}
}
