package rest

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/braintree/manners"
	"github.com/gorilla/mux"
	"github.com/kiasaki/batbelt/http/mm"
)

var logger *log.Logger

func init() {
	logger = log.New(
		os.Stdout,
		fmt.Sprintf("pid:%d ", syscall.Getpid()),
		log.Ldate|log.Lmicroseconds|log.Lshortfile,
	)
}

type Server struct {
	Router      *mux.Router
	AdminRouter *mux.Router
	Filters     mm.Chain
}

func NewServer() Server {
	return Server{
		Router:      mux.NewRouter(),
		AdminRouter: mux.NewRouter(),
		Filters:     mm.New(),
	}
}

func (s *Server) AddFilters(m ...mm.Middleware) {
	s.Filters.Append(m...)
}

// Register in the current server's router all methods handled by
// given endpoint (implementing GET, POST, PUT, DELETE, HEAD)
func (s *Server) Register(endpoint interface{}) {
	RegisterEnpointToRouter(s.Router, endpoint)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Starts accepting http request for the web server and the admin web server
//
// - web: "0.0.0.0:${PORT:-8080}"
// - admin: "127.0.0.1:${PORT:-8081}"
func (s *Server) Run() {
	var wg sync.WaitGroup

	webServer := manners.NewServer()
	adminServer := manners.NewServer()

	go func() {
		wg.Add(1)
		defer wg.Done()
		logger.Println("Web server listening on 0.0.0.0:" + getEnv("PORT", "8080"))
		webServer.ListenAndServe("0.0.0.0:"+getEnv("PORT", "8080"), s.Filters.Then(s.Router))
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		logger.Println("Admin web server listening on 127.0.0.1:" + getEnv("ADMIN_PORT", "8081"))
		adminServer.ListenAndServe("127.0.0.1:"+getEnv("ADMIN_PORT", "8081"), s.AdminRouter)
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// When we get an signal say to both servers to shutdown and
	// wait for the both to finish
	<-signalChan
	webServer.Shutdown <- true
	adminServer.Shutdown <- true

	// Now servers know they need to shutdown just wait till they are done
	wg.Wait()
}
