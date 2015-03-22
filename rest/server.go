package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/braintree/manners"
	"github.com/gorilla/mux"
	"github.com/kiasaki/batbelt/http/mm"
)

type Server struct {
	AppName     string
	Version     string
	Router      *mux.Router
	AdminRouter *mux.Router
	Filters     mm.Chain
	Logger      *log.Logger
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetNotFoundResponse(w)
	})
	return r
}

func newLogger(name, version string) *log.Logger {
	return log.New(
		os.Stdout,
		fmt.Sprintf("app=%s v=%s ", name, version),
		log.Ldate|log.Lmicroseconds|log.Lshortfile,
	)
}

func NewServer(name, version string) Server {
	s := Server{
		AppName:     name,
		Version:     version,
		Router:      newRouter(),
		AdminRouter: newRouter(),
		Filters:     mm.New(),
		Logger:      newLogger(name, version),
	}
	s.AddFilters(s.log)
	return s
}

func (s *Server) AddFilters(m ...mm.Middleware) {
	s.Filters = s.Filters.Append(m...)
}

// Register in the current server's router all methods handled by
// given endpoint (implementing GET, POST, PUT, DELETE, HEAD)
func (s *Server) Register(endpoint interface{}) {
	if e, ok := endpoint.(GET); ok {
		s.Logger.Printf("Registering [%T] at path [%s]\n", endpoint, e.Path())
	}
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
		s.Logger.Println("Web server listening on 0.0.0.0:" + getEnv("PORT", "8080"))
		webServer.ListenAndServe("0.0.0.0:"+getEnv("PORT", "8080"), s.Filters.Then(s.Router))
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		s.Logger.Println("Admin web server listening on 127.0.0.1:" + getEnv("ADMIN_PORT", "8081"))
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

func (s *Server) log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(started time.Time) {
			timing := time.Since(started).Nanoseconds() / 1000.0
			s.Logger.Printf("%s: %s (%dus)\n", r.Method, r.RequestURI, timing)
		}(time.Now())
		handler.ServeHTTP(w, r)
	})
}
