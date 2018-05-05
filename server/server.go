package server

import (
	"github.com/gorilla/mux"

	"net/http"

	doorapi "github.com/tlmiller/garage-door-controller/door/api"
)

type Server struct {
	router *mux.Router
}

func InitServer() *Server {
	server := Server{mux.NewRouter()}
	apiRoute := server.router.PathPrefix("/api")
	doorapi.GetRoutes(apiRoute.Subrouter())

	return &server
}

func (s *Server) Start() {
	http.ListenAndServe(":8787", s.router)
}
