package controller

import (
	internalMiddleware "github.com/adwinugroho/test-chat-multi-schema/controller/middleware"
	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo, userHandler *UserHandler) {
	auth := e.Group("/authentication")
	auth.POST("/login", userHandler.LoginUser)
}

func TenantRoutes(e *echo.Echo, tenantHandler *TenantHandler, userSvc domain.UserService) {
	tenants := e.Group("/tenants")
	tenants.Use(internalMiddleware.AuthenticationMiddleware(userSvc))
	tenants.POST("", tenantHandler.NewTenant)
	tenants.DELETE("/:id", tenantHandler.RemoveTenantByID)
	tenants.PUT("/:id", tenantHandler.UpdateTenantConcurrency)
}

func MessageRoutes(e *echo.Echo, messageHandler *MessageHandler) {
	messages := e.Group("/messages")
	messages.Use(internalMiddleware.ValidateTenantMiddleware())
	messages.POST("/publish", messageHandler.PublishMessage)
	messages.GET("", messageHandler.ListMessages)
}
