package rxk

import (
	"wartracker/ui/wartracker-ui/handler"

	"github.com/go-chi/chi/v5"
)

func RxKRoutes() chi.Router {
	r := chi.NewRouter()

	h := handler.Handler{}
	r.Get("/*", h.Default)

	return r
}
