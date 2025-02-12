package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"wartracker/pkg/vsduel"

	"github.com/go-chi/chi/v5"
)

type VsDuelHandler struct {
}

func (h VsDuelHandler) ListVsDuel(w http.ResponseWriter, r *http.Request) {
	http.Error(w, ErrNotImplemented.Error(), http.StatusNotImplemented)
}

func (h VsDuelHandler) GetVsDuel(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var err error
	var vs vsduel.VsDuel
	err = vs.GetById(chi.URLParam(r, "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, ErrVsDuelNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, vs, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h VsDuelHandler) CreateVsDuel(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var vs vsduel.VsDuel

	err := json.NewDecoder(r.Body).Decode(&vs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = vs.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, vs, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h VsDuelHandler) UpdateVsDuel(w http.ResponseWriter, r *http.Request) {
	http.Error(w, ErrNotImplemented.Error(), http.StatusNotImplemented)
}

func (h VsDuelHandler) AddVsCommanderData(w http.ResponseWriter, r *http.Request) {
	http.Error(w, ErrNotImplemented.Error(), http.StatusNotImplemented)
}

func (h VsDuelHandler) ScanVsDay(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var err error
	var z []byte

	z, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var d vsduel.VsDuel

	d.GetById(chi.URLParam(r, "id"))
	cd, err := d.ScanPointsRanking(z, chi.URLParam(r, "day"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonOutput(w, cd, indent)
}

func VsDuelRoutes() chi.Router {
	r := chi.NewRouter()

	vsDuelHandler := VsDuelHandler{}
	r.Get("/", vsDuelHandler.ListVsDuel)
	r.Post("/", vsDuelHandler.CreateVsDuel)
	r.Post("/scan/{id}/{day}", vsDuelHandler.ScanVsDay)
	r.Get("/{id}", vsDuelHandler.GetVsDuel)
	r.Put("/{id}", vsDuelHandler.UpdateVsDuel)
	r.Put("/{id}/commanderdata", vsDuelHandler.AddVsCommanderData)

	return r
}
