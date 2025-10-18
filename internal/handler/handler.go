package handler

import (
	"L0-wb/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service service.Service
}

func NewHandler(service service.Service) Handler {
	return &UserHandler{service: service}
}

func (h *UserHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

var ErrNotFound = errors.New("order not found")

// refactor cpntext
func (h *UserHandler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID, ok := vars["uid"]
	if !ok || orderUID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"msg":    "Missing order UID",
		})
		return
	}

	ctx := r.Context() //refactor
	order, err := h.service.GetOrderByUID(ctx, orderUID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]interface{}{
				"status": "error",
				"msg":    "Order not found",
			})
			log.Printf("order with UID %s not found", orderUID)
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status": "error",
			"msg":    "Internal server error",
		})
		log.Printf("failed to retrieve order with UID %s: %v", orderUID, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"data":   order,
	})
	log.Printf("order with UID %s retrieved successfully", orderUID)
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
