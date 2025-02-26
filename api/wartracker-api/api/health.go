package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthHandler struct {
}

func (h HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK\n"))
}

func HealthRoutes() chi.Router {
	r := chi.NewRouter()

	healthHandler := HealthHandler{}
	r.Get("/", healthHandler.GetHealth)

	return r
}
