package common

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type OsNameKeyType struct{}

func GetOsName(req *http.Request) string {
	// リクエストヘッダーからUser-Agentを取得
	userAgent := req.Header.Get("User-Agent")
	// useragentを使ってOSを取得
	OsName := useragent.Parse(userAgent).OS
	return OsName
}

// 取得したOsNameをcontextに格納する
func SetOsName(ctx context.Context, OsName string) context.Context {
	return context.WithValue(ctx, OsNameKeyType{}, OsName)
}
