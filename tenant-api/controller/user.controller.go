package controller

import (
	"net/http"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(svc domain.UserService) UserHandler {
	return UserHandler{service: svc}
}

func (h *UserHandler) LoginUser(c echo.Context) error {
	var req model.LoginUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, "Bad request"))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	resp, err := h.service.LoginUser(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).
		SetMessage("Successfully logged in").
		SetData(resp))
}
