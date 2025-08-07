package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adwinugroho/test-chat-multi-schema/config"
	"github.com/adwinugroho/test-chat-multi-schema/controller"
	internalMiddleware "github.com/adwinugroho/test-chat-multi-schema/controller/middleware"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/server"
	"github.com/adwinugroho/test-chat-multi-schema/repository"
	"github.com/adwinugroho/test-chat-multi-schema/service"
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
	logger.LogInfo("Starting application initialization...")

	config.LoadConfig()

	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second) // Reduced timeout
	defer cancel()

	// Initialize database
	dbHandler, err := config.InitConnectDB(ctx, config.PostgreSQLConfig.Database.URL)
	if err != nil {
		logger.LogFatal("Failed to connect to database:" + err.Error())
	}
	defer dbHandler.CloseAllConnection()

	// Initialize RabbitMQ
	err = config.InitRabbitMQConnection(config.RabbitMQConfig.RabbitMQ.URL)
	if err != nil {
		logger.LogFatal("Failed to connect to rabbitMQ:" + err.Error())
	}
	defer func() {
		if config.RabbitConn != nil {
			config.RabbitConn.Close()
			logger.LogInfo("RabbitMQ connection closed")
		}
	}()

	e := echo.New()
	e.Validator = &internalMiddleware.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self' https: 'unsafe-inline'")
			return next(c)
		}
	})

	userRepository := repository.NewUserRepository(dbHandler.DB)
	userService := service.NewUserService(userRepository)
	authHandler := controller.NewUserHandler(userService)

	tenantRepository := repository.NewTenantRepository(dbHandler.DB)
	messageRepository := repository.NewMessageRepository(dbHandler.DB)
	tenantService := service.NewTenantService(tenantRepository)
	subscriberService := service.NewListenSubscriber(messageRepository)
	tenantManager := service.NewTenantManager(subscriberService)
	tenantHandler := controller.NewTenantHandler(tenantService, tenantManager)

	controller.UserRoutes(e, authHandler)
	controller.TenantRoutes(e, tenantHandler, userService)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":   "OK",
			"database": "connected",
			"rabbitmq": "connected",
		})
	})

	logger.LogInfo("Application fully initialized")
	logger.LogInfo("Server starting on port:" + config.AppConfig.Port)
	logger.LogInfo("Application URL:" + config.AppConfig.AppURL)

	server.StartServerWithGracefulShutdown(e, ctx, dbHandler)
}
