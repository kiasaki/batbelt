package middlewares

import "net/http"

func BasicAuth(user string, pass string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if username, password, ok := r.BasicAuth(); ok {
				if username != user || password != pass {
					http.Error(w, "401 Bad authorization", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "401 Missing authorization", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
