package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adwinugroho/test-chat-multi-schema/controller"
	internalMiddleware "github.com/adwinugroho/test-chat-multi-schema/controller/middleware"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTenantManager struct {
	mock.Mock
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger.InitLogger()
}

func (m *mockTenantManager) RestartConsumer(ctx context.Context, tenantID string, workers int) error {
	args := m.Called(ctx, tenantID, workers)
	return args.Error(0)
}

func (m *mockTenantManager) StartConsumer(ctx context.Context, tenantID string, workers int) error {
	args := m.Called(ctx, tenantID, workers)
	return args.Error(0)
}

func (m *mockTenantManager) StopConsumer(tenantID string) error {
	args := m.Called(tenantID)
	return args.Error(0)
}

func (m *mockTenantManager) StopAllConsumers() {}

func setupEchoTest() *echo.Echo {
	e := echo.New()
	e.Validator = &internalMiddleware.CustomValidator{Validator: validator.New()}
	return e
}

func TestUpdateTenantConcurrency_Success(t *testing.T) {
	e := setupEchoTest()

	mockManager := new(mockTenantManager)
	h := controller.NewTenantHandler(nil, mockManager)

	reqBody := model.UpdateTenantConcurrencyRequest{Workers: 5}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/tenants/abc", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues("abc")

	mockManager.On("RestartConsumer", mock.Anything, "abc", 5).Return(nil)

	err := h.UpdateTenantConcurrency(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockManager.AssertExpectations(t)
}

func TestUpdateTenantConcurrency_ErrorRestart(t *testing.T) {
	e := setupEchoTest()

	mockManager := new(mockTenantManager)
	h := controller.NewTenantHandler(nil, mockManager)

	reqBody := model.UpdateTenantConcurrencyRequest{Workers: 2}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/tenants/abc", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues("abc")

	mockManager.On("RestartConsumer", mock.Anything, "abc", 2).Return(errors.New("restart failed"))

	err := h.UpdateTenantConcurrency(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockManager.AssertExpectations(t)
}
