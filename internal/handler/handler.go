package handler

import (
	"L0-wb/internal/models"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error)
}

type UserHandler struct {
	service Service
}

func NewHandler(service Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

var ErrNotFound = errors.New("order not found")

func (h *UserHandler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID, ok := vars["uid"]
	if !ok || orderUID == "" {
		http.Error(w, "Missing order UID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	order, err := h.service.GetOrderByUID(ctx, orderUID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
			log.Printf("order with UID %s not found", orderUID)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to retrieve order with UID %s: %v", orderUID, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
	log.Printf("order with UID %s retrieved successfully", orderUID)
}
