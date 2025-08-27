package handler_test

import (
	"L0-wb/internal/handler"
	"L0-wb/internal/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	order *models.Order
	err   error
}

func (m *mockService) GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	return m.order, m.err
}

func TestServeIndex(t *testing.T) {
	// создаём временный index.html
	testFile := "./web/index.html"
	_ = os.MkdirAll("./web", 0755)
	_ = os.WriteFile(testFile, []byte("hello index"), 0644)
	defer os.Remove(testFile)

	h := handler.NewHandler(&mockService{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	h.ServeIndex(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "hello index", rr.Body.String())
}

func TestGetOrderByUID_Success(t *testing.T) {
	expectedOrder := &models.Order{OrderUID: "12345"}

	h := handler.NewHandler(&mockService{order: expectedOrder})
	req := httptest.NewRequest(http.MethodGet, "/order/12345", nil)
	rr := httptest.NewRecorder()

	// эмулируем gorilla/mux Vars
	req = mux.SetURLVars(req, map[string]string{"uid": "12345"})

	h.GetOrderByUID(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var got models.Order
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	require.Equal(t, expectedOrder.OrderUID, got.OrderUID)
}

func TestGetOrderByUID_NotFound(t *testing.T) {
	h := handler.NewHandler(&mockService{order: nil, err: handler.ErrNotFound})
	req := httptest.NewRequest(http.MethodGet, "/order/00000", nil)
	rr := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"uid": "00000"})

	h.GetOrderByUID(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Contains(t, rr.Body.String(), "Order not found")
}

func TestGetOrderByUID_InternalError(t *testing.T) {
	h := handler.NewHandler(&mockService{order: nil, err: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/order/99999", nil)
	rr := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"uid": "99999"})

	h.GetOrderByUID(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Contains(t, rr.Body.String(), "Internal server error")
}
