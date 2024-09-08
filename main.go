package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"golang.org/x/sync/errgroup"
)

// 環境変数から取得した値を使ってサーバーを起動する
var (
	// Basic Auth用のユーザー名とパスワードを環境変数から取得
	username = os.Getenv("BASIC_AUTH_USER_ID")
	password = os.Getenv("BASIC_AUTH_PASSWORD")
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする.
	mux := router.NewRouter(todoDB, username, password)

	// シグナルを受け取るためのコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	// TODO: サーバーをlistenする
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// errgroupを作成
	g, ctx := errgroup.WithContext(ctx)

	// シグナルを受け取り、サーバーをシャットダウンするゴルーチンをerrgroupで実行
	g.Go(func() error {
		// シグナルを受け取るまで待機
		<-ctx.Done()

		// 5秒のタイムアウト付きコンテキストを作成
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// サーバーをシャットダウン(新しい接続の受け付けを停止し、contextがキャンセルされたら終了する)
		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	// メインの処理としてサーバーを起動し、正常に終了しない場合はエラーを返す
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	// サーバーがシャットダウンされるゴルーチンが終了するまで待機
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
