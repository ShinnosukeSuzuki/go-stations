package handler

import (
	"net/http"
)

// A DoPanicHandler implements do panic endpoint.
type DoPanicHandler struct{}

// NewDoPanicHandler returns DoPanicHandler based http.Handler.
func NewDoPanicHandler() *DoPanicHandler {
	return &DoPanicHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *DoPanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ivoke panic
	panic("do-panic")
}
