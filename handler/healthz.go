package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 5秒後にレスポンスを返す
	time.Sleep(10 * time.Second)
	response := &model.HealthzResponse{
		Message: "OK",
	}
	json.NewEncoder(w).Encode(response)
}
