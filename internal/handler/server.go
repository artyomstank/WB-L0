package handler

import (
	"fmt"
	"net/http"

	"L0-wb/config"

	"github.com/gorilla/mux"
)

func NewServer(cfg *config.Config, h *UserHandler) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/", h.ServeIndex).Methods(http.MethodGet)
	router.HandleFunc("/order/{uid}", h.GetOrderByUID).Methods(http.MethodGet)

	addrStr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	return &http.Server{
		Addr:         addrStr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}
}
