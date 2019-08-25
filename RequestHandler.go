package pipeflow

import (
	"net/http"
	"strings"
)

// HTTPMethod is enum of http methods
type HTTPMethod int

const (
	// HTTPGet GET
	HTTPGet = iota
	// HTTPHead HEAD
	HTTPHead
	// HTTPPost POST
	HTTPPost
	// HTTPPut PUT
	HTTPPut
	// HTTPDelete DELETE
	HTTPDelete
	// HTTPConnect CONNECT
	HTTPConnect
	// HTTPOptions OPTIONS
	HTTPOptions
	// HTTPTrace TRACE
	HTTPTrace
)

// RequestHandler is used to register a request handler
type RequestHandler struct {
	Route   *Route
	Methods map[HTTPMethod]bool
	Handle  func(ctx HTTPContext)
}

// Conflict checks handler's path equals to other's and HTTP methods have intersection
func (h *RequestHandler) Conflict(other *RequestHandler) bool {
	if h.Route.Equals(other.Route) {
		return h.HasInterMethod(other)
	}

	return false
}

// HasInterMethod checks whether http methods has intersection
func (h *RequestHandler) HasInterMethod(other *RequestHandler) bool {
	for k := range h.Methods {
		if _, ok := other.Methods[k]; ok {
			return true
		}
	}

	return false
}

// Match check there is any request handler.
func (h *RequestHandler) Match(request *http.Request) bool {
	path := request.URL.Path
	method := request.Method

	if !h.Route.PathReg.MatchString(path) {
		return false
	}

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}
	httpMethods := []HTTPMethod{HTTPGet, HTTPHead, HTTPPost, HTTPPut, HTTPDelete, HTTPConnect, HTTPOptions, HTTPTrace}

	method = strings.ToUpper(method)
	httpMethod := -1
	for i, v := range methods {
		if v == method {
			httpMethod = i
			break
		}
	}

	if -1 != httpMethod {
		hasInter := h.HasInterMethod(&RequestHandler{Methods: map[HTTPMethod]bool{httpMethods[httpMethod]: true}})
		if !hasInter {
			return false
		}
	}

	if e := request.ParseForm(); e != nil {
		return false
	}
	for k := range h.Route.Params {
		if _, ok := request.Form[k]; !ok {
			return false
		}
	}

	return true
}
