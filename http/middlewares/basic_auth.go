package middlewares

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"
)

// Authenticate func from: https://github.com/goji/httpauth/blob/master/basic_auth.go
func authenticate(user, pass string, r *http.Request) bool {
	const basicScheme string = "Basic "

	// Confirm the request is sending Basic Authentication credentials.
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, basicScheme) {
		return false
	}

	// Get the plain-text username and password from the request
	// The first six characters are skipped - e.g. "Basic ".
	str, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
	if err != nil {
		return false
	}

	// Split on the first ":" character only, with any subsequent colons assumed to be part
	// of the password. Note that the RFC2617 standard does not place any limitations on
	// allowable characters in the password.
	creds := bytes.SplitN(str, []byte(":"), 2)

	if len(creds) != 2 {
		return false
	}

	// Equalize lengths of supplied and required credentials
	// by hashing them
	givenUser := sha256.Sum256(creds[0])
	givenPass := sha256.Sum256(creds[1])
	requiredUser := sha256.Sum256([]byte(user))
	requiredPass := sha256.Sum256([]byte(pass))

	// Compare the supplied credentials to those set in our options
	if subtle.ConstantTimeCompare(givenUser[:], requiredUser[:]) == 1 &&
		subtle.ConstantTimeCompare(givenPass[:], requiredPass[:]) == 1 {
		return true
	}

	return false
}

func BasicAuth(user string, pass string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !authenticate(user, pass, r) {
				w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				http.Error(w, "401 Bad authorization", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
