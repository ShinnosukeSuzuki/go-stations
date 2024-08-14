package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/common"
)

// ログに出力する構造体を定義
type AccessLogging struct {
	Timestamp time.Time
	Latency   int64
	Path      string
	OS        string
}

// アクセス日時、リクエストパス、処理時間を出力するミドルウェア
func AccessLoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// OsNameを取得
		osName := common.GetOsName(r)
		// contextに格納
		ctx := common.SetOsName(r.Context(), osName)
		r = r.WithContext(ctx)

		// handlerのアクセス時刻を取得
		start := time.Now()
		// ハンドラに渡す
		h.ServeHTTP(w, r)
		// handlerの処理時間を取得
		duration := time.Since(start)
		// AccessLoggingにアクセス日時、処理時間、リクエストパスを格納
		al := AccessLogging{
			Timestamp: start,
			Latency:   duration.Milliseconds(),
			Path:      r.URL.Path,
			OS:        osName,
		}

		fmt.Printf("%+v\n", al)
	})
}
