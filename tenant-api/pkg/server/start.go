package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adwinugroho/test-chat-multi-schema/config"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/labstack/echo/v4"
)

func StartServerWithGracefulShutdown(e *echo.Echo, ctx context.Context, db *config.PostgresDB) {
	idleConnsClosed := make(chan struct{})

	logger.LogInfo(fmt.Sprintf("Server is starting on port %s", config.AppConfig.Port))

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
		sig := <-sigint

		logger.LogInfo(fmt.Sprintf("Received shutdown signal: %s", sig))
		logger.LogInfo("Initiating graceful shutdown sequence")

		logger.LogInfo("Step 1: Closing database connections")
		if db != nil {
			db.CloseAllConnection()
			logger.LogInfo("All database connections have been closed successfully")
		} else {
			logger.LogInfo("No database connection to close")
		}

		logger.LogInfo("Step 2: Shutting down HTTP server")
		if err := e.Shutdown(ctx); err != nil {
			logger.LogError("HTTP server shutdown error: " + err.Error())
		} else {
			logger.LogInfo("HTTP server shutdown completed successfully")
		}

		logger.LogInfo("Graceful shutdown completed")
		close(idleConnsClosed)
	}()

	if err := e.Start(fmt.Sprintf(":%s", config.AppConfig.Port)); err != nil {
		if err.Error() != "http: Server closed" {
			logger.LogError("Server error: " + err.Error())
		}
	}

	<-idleConnsClosed
	logger.LogInfo("Server process exiting")
}

func StartServer(e *echo.Echo) {
	if err := e.Start(fmt.Sprintf(":%s", config.AppConfig.Port)); err != nil {
		logger.LogError("Oops... Server is not running! Reason: " + err.Error())
	}
}
