package controller

import (
	"net/http"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TenantHandler struct {
	service domain.TenantService
}

func NewTenantHandler(svc domain.TenantService) TenantHandler {
	return TenantHandler{service: svc}
}

func (h *TenantHandler) NewTenant(c echo.Context) error {
	var req model.NewTenantRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, "Bad request"))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	payload := &domain.Tenant{
		TenantID:   uuid.New().String(),
		TenantName: req.TenantName,
	}

	err := h.service.NewTenant(c.Request().Context(), payload)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).
		SetMessage("Successfully created new tenant"))
}

func (h *TenantHandler) RemoveTenantByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, "Bad request"))
	}
	err := h.service.RemoveTenantByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).
		SetMessage("Successfully remove tenant"))
}
