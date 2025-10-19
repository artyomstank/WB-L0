package handler

import (
	"L0-wb/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service service.Service
}

func NewHandler(service service.Service) Handler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

var ErrNotFound = errors.New("order not found")

func (h *UserHandler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID, ok := vars["uid"]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"status": "error",
			"msg":    "Missing order UID",
		})
		return
	}

	// Убираем пробелы с начала и конца
	orderUID = strings.TrimSpace(orderUID)
	if orderUID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"msg":    "Order UID cannot be empty",
		})
		return
	}

	ctx := r.Context()
	order, err := h.service.GetOrderByUID(ctx, orderUID)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Internal server error"

		if errors.Is(err, service.ErrNotFound) {
			status = http.StatusNotFound
			msg = "Order not found"
		}

		writeJSON(w, status, map[string]interface{}{
			"status": "error",
			"msg":    msg,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"data":   order,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
