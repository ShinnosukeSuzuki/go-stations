package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	// /healthzの時にHealthzHandlerを呼び出す
	mux.Handle("/healthz", handler.NewHealthzHandler())

	//todoDBを使ってserviceを作成
	todoService := service.NewTODOService(todoDB)
	mux.Handle("/todos", handler.NewTODOHandler(todoService))

	return mux
}
