package muxserver

import (
	"net/http"
)

type CustomerServerMux struct {
	http.ServeMux
	middleware []func(next http.Handler) http.Handler
}

func (c *CustomerServerMux) RegisterMiddleware(next func(next http.Handler) http.Handler) {
	c.middleware = append(c.middleware, next)
}

func (c *CustomerServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var current http.Handler = &c.ServeMux

	for _, next := range c.middleware {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}
