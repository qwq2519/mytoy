package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	do "github.com/samber/do/v2"

	"mytoy/bootstrap"
	"mytoy/config"
	"mytoy/controller"
)

func main() {
	// 初始化依赖注入容器（包括数据库等）
	injector := bootstrap.NewContainer()

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	cfgManager, err := do.Invoke[*config.Manager](injector)
	if err != nil {
		logger.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	snapshot := cfgManager.Snapshot()

	if snapshot.Server.Mode == "" {
		logger.Error("server.mode is required")
		os.Exit(1)
	}
	gin.SetMode(snapshot.Server.Mode)

	// 初始化路由
	router := controller.NewRouter(injector)

	if snapshot.Server.Port <= 0 {
		logger.Error("server.port is required and must be > 0", "port", snapshot.Server.Port)
		os.Exit(1)
	}
	addr := fmt.Sprintf(":%d", snapshot.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	logger.Info("starting server", "addr", addr)

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if err != nil {
			logger.Error("failed to start server", "err", err)
			os.Exit(1)
		}
		return
	case <-sigCtx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "err", err)
	}
	if err := <-serverErr; err != nil {
		logger.Error("server exited with error", "err", err)
	}

	if report := injector.ShutdownWithContext(shutdownCtx); report != nil && len(report.Errors) > 0 {
		logger.Error("di shutdown finished with errors", "err", report.Error())
	}
}
