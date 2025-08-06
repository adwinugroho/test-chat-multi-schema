package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adwinugroho/test-chat-multi-schema/config"
	controllerUser "github.com/adwinugroho/test-chat-multi-schema/controller"
	internalMiddleware "github.com/adwinugroho/test-chat-multi-schema/controller/middleware"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	repoUser "github.com/adwinugroho/test-chat-multi-schema/repository"
	serviceUser "github.com/adwinugroho/test-chat-multi-schema/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger.InitLogger()

	config.LoadConfig()
}

func main() {
	logger.LogInfo("Starting application...")
	logger.LogInfo("Application started on port:" + config.AppConfig.Port)
	logger.LogInfo("Application URL:" + config.AppConfig.AppURL)

	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 60*time.Second)
	defer cancel()

	dbHandler, err := config.InitConnectDB(ctx, config.PostgreSQLConfig.Database.URL)
	if err != nil {
		logger.LogFatal("Failed to connect to database:" + err.Error())
	}

	var e = echo.New()
	e.Validator = &internalMiddleware.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Recover())

	userRepository := repoUser.NewUserRepository(dbHandler.DB)
	userService := serviceUser.NewUserService(userRepository)
	authHandler := controllerUser.NewUserHandler(userService)

	controllerUser.UserRoutes(e, authHandler, userService)

	e.Use(middleware.Secure())

	// Add security headers middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Security headers
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self' https: 'unsafe-inline'")

			return next(c)
		}
	})

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.AppConfig.Port)))
}
