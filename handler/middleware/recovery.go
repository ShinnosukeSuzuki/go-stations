package middleware

import (
	"log"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// severHTTP内でpanicが発生した場合でもdeferは実行される
		// recoverはdeferの中でのみ使用可能
		defer func() {
			if err := recover(); err != nil {
				// panic理由とURLをログに出力
				log.Printf("panic: %v, URL: %s", err, r.URL.String())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
