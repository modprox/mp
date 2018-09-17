package webutil

import (
	"fmt"
	"net/http"
)

// A Middleware is used to execute intermediate Handlers
// in response to a request.
type Middleware func(http.Handler) http.Handler

//  Chain recursively chains middleware together.
func Chain(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](Chain(h, m[1:cap(m)]...))
}

const (
	HeaderAPIKey = "X-modprox-api-key"
)

// KeyGuard creates a Middleware which protects access to a handler
// by first checking for the X-modprox-api-key header being set. If
// at least one of the values for the header matches one of the keys
// configured for the KeyGuard, the handler is executed for the request.
// Otherwise, a StatusForbidden response is returned.
func KeyGuard(keys []string) Middleware {
	allowedKeys := make(map[string]bool)
	for _, key := range keys {
		allowedKeys[key] = true
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// first check that the header is set
			key := r.Header.Get(HeaderAPIKey)
			if key == "" {
				msg := fmt.Sprintf("header %s is not set in request", HeaderAPIKey)
				http.Error(w, msg, http.StatusForbidden)
				return
			}

			// check if the given key is allowable
			if allowedKeys[key] {
				// found a good key, execute the
				// protected handler for the request
				h.ServeHTTP(w, r)
				return
			}

			// no good key was provided, respond with an error
			msg := fmt.Sprintf("header %s contains no valid keys", HeaderAPIKey)
			http.Error(w, msg, http.StatusForbidden)
			return
		})
	}
}
