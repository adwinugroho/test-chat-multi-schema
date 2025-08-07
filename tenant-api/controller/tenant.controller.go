package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/adwinugroho/test-chat-multi-schema/config"
	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/adwinugroho/test-chat-multi-schema/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TenantHandler struct {
	service       domain.TenantService
	tenantManager *service.TenantManager
}

func NewTenantHandler(svc domain.TenantService, tm *service.TenantManager) *TenantHandler {
	return &TenantHandler{service: svc, tenantManager: tm}
}

func (h *TenantHandler) NewTenant(c echo.Context) error {
	var req model.NewTenantRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, "Bad request"))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userSession, ok := c.Get("user").(*domain.User)
	if !ok {
		return c.JSON(http.StatusForbidden, model.NewError(model.ErrorUnauthorized, "Access denied"))
	}

	payload := &domain.Tenant{
		TenantID:   uuid.New().String(),
		TenantName: req.TenantName,
		UserID:     userSession.UserID,
	}

	err := h.service.NewTenant(c.Request().Context(), payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = h.service.CreateTenantPartition(c.Request().Context(), payload.TenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	defaultWorkers, err := strconv.Atoi(config.AppConfig.Workers)
	if err != nil {
		logger.LogError("Default worker is invalid number:" + err.Error())
		return c.JSON(http.StatusInternalServerError, model.NewError(model.ErrorGeneral, "Internal server error"))
	}

	// replace context
	consumerCtx := context.Background()
	err = h.tenantManager.StartConsumer(consumerCtx, payload.TenantID, defaultWorkers)
	if err != nil {
		logger.LogError("Error while start consumer:" + err.Error())
		return c.JSON(http.StatusInternalServerError, model.NewError(model.ErrorGeneral, "Internal server error"))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).
		SetMessage("Successfully created new tenant, tenant consumer has started...").
		SetData(payload))
}

func (h *TenantHandler) RemoveTenantByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, "Bad request"))
	}

	err := h.tenantManager.StopConsumer(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewError(model.ErrorGeneral, err.Error()))
	}

	err = h.service.RemoveTenantByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).
		SetMessage("Successfully remove tenant"))
}

func (h *TenantHandler) UpdateTenantConcurrency(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, "Bad request"))
	}

	var req model.UpdateTenantConcurrencyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorInvalidRequest, "Invalid request body"))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, err.Error()))
	}

	err := h.tenantManager.RestartConsumer(c.Request().Context(), id, req.Workers)
	if err != nil {
		logger.LogError("failed to restart consumer: " + err.Error())
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).
		SetMessage("Successfully updated tenant concurrency"))
}
