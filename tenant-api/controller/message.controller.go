package controller

import (
	"net/http"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	service        domain.MessageService
	publishService domain.PublisherService
}

func NewMessageHandler(svc domain.MessageService, publishSvc domain.PublisherService) *MessageHandler {
	return &MessageHandler{
		service:        svc,
		publishService: publishSvc,
	}
}

func (h *MessageHandler) PublishMessage(c echo.Context) error {
	var req model.PublishRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewError(model.ErrorBadRequest, "Bad request"))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	tenantID := c.Get("tenant_id").(string)
	h.service.PublishMessage(c.Request().Context(), tenantID, &req)

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).
		SetMessage("Successfully publish message"))
}
