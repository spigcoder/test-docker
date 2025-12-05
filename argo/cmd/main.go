// 程序入口
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/swanhubx/swanlab-helper/argo/pkg/config"
	"github.com/swanhubx/swanlab-helper/argo/pkg/logger"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/swanhubx/swanlab-helper/argo/internal/router"
)

func main() {
	// 使用 UTC 时区
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	err := config.Init("configs", "config", "ARGO")
	if err != nil {
		panic(err)
	}
	logger.Init("info")

	// 初始化数据库实例
	if !viper.IsSet("database.url") {
		slog.Error("Database URL not set")
	}
	db, err := gorm.Open(mysql.Open(viper.GetString("database.url")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	env := router.Env{
		DB: db,
	}
	// 初始化 router
	r := router.NewRouter(env)
	addr := viper.GetString("addr")
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ListenAndServe: ", "error", err)
			panic(err)
		}
	}()
	// 监听推出信号，优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", "error: ", err)
	}
	slog.Info("Server exiting")
}
