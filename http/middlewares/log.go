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
		start := time.Now()

		defer func() {
			end := time.Now()
			duration := end.Sub(start)
			if duration < time.Millisecond {
				log.Printf("%s %s %s (%dμs)", r.RemoteAddr, r.Method, r.URL, duration/time.Microsecond)
			} else {
				log.Printf("%s %s %s (%dms)", r.RemoteAddr, r.Method, r.URL, duration/time.Millisecond)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
