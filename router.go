package http

import (
	"net/http"
	"strings"
)

type HandleFunc func(r *Request, w *ResponseWriter)

type Router struct {
	routes map[string]HandleFunc
}

func (r *Router) Handle(rq *Request, w *ResponseWriter) {
	path := strings.TrimSuffix(rq.URL, "/")
	if h, contains := r.routes[path]; contains {
		h(rq, w)
		return
	}

	_ = w.WriteBody(http.StatusNotFound, []byte("not found"))
}

func (r *Router) HandleFunc(path string, handler HandleFunc) {
	if r.routes == nil {
		r.routes = map[string]HandleFunc{path: handler}
	}

	r.routes[path] = handler
}
