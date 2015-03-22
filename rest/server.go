package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiasaki/batbelt/http/mm"
)

type Server struct {
	Router  *mux.Router
	Filters mm.Chain
}

func NewServer() Server {
	return Server{
		Router:  mux.NewRouter(),
		Filters: mm.New(),
	}
}

// Helper to fetch vars for current request from context
func (s *Server) Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func (s *Server) AddFilters(m ...mm.Middleware) {
	s.Filters.Append(m...)
}

// Register in the current server's router all methods handled by
// given endpoint (implementing GET, POST, PUT, DELETE, HEAD)
func (s *Server) Register(endpoint interface{}) {
	RegisterEnpointToRouter(s.Router, endpoint)
}

// Start listening in configured ports
func (s *Server) Run() {

}
