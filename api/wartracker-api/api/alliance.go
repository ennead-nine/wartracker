package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"wartracker/pkg/alliance"

	"github.com/go-chi/chi/v5"
)

type AllianceHandler struct {
}

func (h AllianceHandler) List(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	as, err := alliance.List()
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, ErrNoAliiances.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, as, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h AllianceHandler) Get(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")
	latest := GetQueryBool(r, "latest")

	var a alliance.Alliance
	a.Id = chi.URLParam(r, "id")

	err := a.Get()
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, ErrAliianceNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if latest {
		err = a.GetLatestData()
	} else {
		err = a.GetData()
	}
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, a, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h AllianceHandler) Create(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var a alliance.Alliance

	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var date string
	for k := range a.DataMap {
		date = k
	}

	err = a.AddData(date, a.DataMap[date])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, a, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h AllianceHandler) Update(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var a alliance.Alliance

	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, a, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h AllianceHandler) AddData(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var a alliance.Alliance
	a.Id = chi.URLParam(r, "id")
	err := a.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var d alliance.Data
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.AddData(d.Date, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, a, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AllianceRoutes() chi.Router {
	r := chi.NewRouter()

	h := AllianceHandler{}
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Put("/{id}/data", h.AddData)

	return r
}
