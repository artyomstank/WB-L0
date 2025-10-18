package handler

import "net/http"

type Handler interface {
	ServeIndex(w http.ResponseWriter, r *http.Request)
	GetOrderByUID(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
}
