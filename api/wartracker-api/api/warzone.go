package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"wartracker/pkg/warzone"

	"github.com/go-chi/chi/v5"
)

type WarzoneHandler struct {
}

func (h WarzoneHandler) List(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	zs, err := warzone.List()
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, ErrNoAliiances.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, zs, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h WarzoneHandler) Create(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var z warzone.Warzone

	s := chi.URLParam(r, "server")
	server, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = z.Create(server)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, z, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h WarzoneHandler) Get(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var z warzone.Warzone
	z.Id = chi.URLParam(r, "id")
	err := z.Get()
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, z, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h WarzoneHandler) GetByServer(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var z warzone.Warzone
	s := chi.URLParam(r, "server")
	server, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = z.GetByServer(server)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, z, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func WarzoneRoutes() chi.Router {
	r := chi.NewRouter()

	h := WarzoneHandler{}
	r.Post("/{server}", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Get("/s/{server}", h.GetByServer)

	return r
}
