package mm

import "net/http"

type Middleware func(http.Handler) http.Handler
type Chain []Middleware

func New(middlewares ...Middleware) Chain {
	c := make([]Middleware, 0)
	return append(c, middlewares...)
}

func (c Chain) Append(middlewares ...Middleware) Chain {
	ms := make([]Middleware, len(c)+len(middlewares))
	copy(ms, c)
	copy(ms[len(c):], middlewares)
	return ms
}

func (c Chain) Then(handler http.Handler) http.Handler {
	final := handler
	for i := len(c) - 1; i >= 0; i-- {
		final = c[i](final)
	}
	return final
}

func (c Chain) ThenFunc(handlerFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	return c.Then(http.HandlerFunc(handlerFunc))
}
