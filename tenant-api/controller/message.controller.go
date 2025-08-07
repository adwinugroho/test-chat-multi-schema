package controller

import (
	"context"
	"net/http"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	service domain.MessageService
}

func NewMessageHandler(svc domain.MessageService) *MessageHandler {
	return &MessageHandler{
		service: svc,
	}
}

func (h *MessageHandler) ListMessages(c echo.Context) error {
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok {
		return c.JSON(http.StatusForbidden, model.NewError(model.ErrorUnauthorized, "Access denied"))
	}

	cursor := c.QueryParam("cursor")
	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	mapQParam := map[string]string{
		"limit":  limitStr,
		"cursor": cursor,
	}

	messages, nextCursor, err := h.service.GetMessages(c.Request().Context(), tenantID, mapQParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewError(model.ErrorGeneral, "Internal server error"))
	}

	if len(messages) == 0 {
		c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(messages).SetMessage("Data not found"))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetListWithCursor(messages, nextCursor))
}

func (h *MessageHandler) PublishMessage(c echo.Context) error {
	var req model.PublishRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, "Bad request"))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok {
		return c.JSON(http.StatusForbidden, model.NewError(model.ErrorUnauthorized, "Access denied"))
	}

	publisherCtx := context.Background()
	err := h.service.PublishMessage(publisherCtx, tenantID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).
		SetMessage("Successfully publish message"))
}
