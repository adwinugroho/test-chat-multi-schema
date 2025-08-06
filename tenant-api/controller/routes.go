package controller

import (
	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo, userHandler UserHandler, userSvc domain.UserService) {
	auth := e.Group("/authentication")
	auth.POST("/login", userHandler.LoginUser)
}
