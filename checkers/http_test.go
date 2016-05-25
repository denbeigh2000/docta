package checkers

import (
	_ "testing"

	"net/http"
)

type simpleHTTPHandler struct {
	f http.HandlerFunc
}

func (h simpleHTTPHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.f(rw, r)
}

func BuildHandler(code int, responseText string) http.Handler {
	return simpleHTTPHandler{func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(code)
		rw.Write([]byte(responseText))
	}}
}
