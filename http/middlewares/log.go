package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func LogWithTiming(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(started time.Time) {
			timing := time.Since(started).Nanoseconds() / 1000.0
			log.Printf("%s: %s (%dus)\n", r.Method, r.RequestURI, timing)
		}(time.Now())
		handler.ServeHTTP(w, r)
	})
}
