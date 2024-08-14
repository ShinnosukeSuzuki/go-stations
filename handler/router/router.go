package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB, username, password string) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	// /healthzの時にHealthzHandlerを呼び出す
	mux.Handle("/healthz", handler.NewHealthzHandler())

	// do-panicの時にmiddlewareのRecoveryを通してDoPanicHandlerを呼び出す
	mux.Handle("/do-panic", middleware.Recovery(handler.NewDoPanicHandler()))

	//todoDBを使ってserviceを作成
	todoService := service.NewTODOService(todoDB)
	mux.Handle("/todos", middleware.BasicAuthMiddleware(middleware.AccessLoggingMiddleware(handler.NewTODOHandler(todoService)), username, password))

	return mux
}
