package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	// /healthzの時にHealthzHandlerを呼び出す
	mux.Handle("/healthz", handler.NewHealthzHandler())
	return mux
}
