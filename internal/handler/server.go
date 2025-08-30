package handler

import (
	"fmt"
	"net/http"

	"L0-wb/config"

	"github.com/gorilla/mux"
)

func NewServer(cfg *config.Config, h *UserHandler) *http.Server {
	router := mux.NewRouter()
	//Middleware для CORS
	router.Use(corsMiddleware)
	// API ручка
	router.HandleFunc("/order/{uid}", h.GetOrderByUID).Methods(http.MethodGet)

	//Статические файлы из папки ./web
	fs := http.FileServer(http.Dir("./web"))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	addrStr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	return &http.Server{
		Addr:         addrStr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}
}

// Добавляем заголовки для кросс-доменных запросов с corsMiddleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // все источники
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		//preflight запрос
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
