package middleware

import (
	"net/http"
)

func BasicAuthMiddleware(h http.Handler, username, password string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Basic認証のユーザー名とパスワードを取得
		user, pass, ok := r.BasicAuth()

		// Basic認証情報がない(または空で送られくる)、またはユーザー名とパスワードが一致しない場合は401 Unauthorizedを返す
		if !ok || user == "" || pass == "" || user != username || pass != password {
			// WWW-Authenticate ヘッダーを設定
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			// 401 Unauthorized ステータスコードを設定
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
