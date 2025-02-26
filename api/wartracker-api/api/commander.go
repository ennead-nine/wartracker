package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"wartracker/pkg/commander"

	"github.com/go-chi/chi/v5"
)

type CommanderHandler struct {
}

func (h CommanderHandler) List(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	as, err := commander.List()
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

func (h CommanderHandler) Get(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var err error
	var c commander.Commander
	c.Id = chi.URLParam(r, "id")
	err = c.Get()
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, c, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h CommanderHandler) Create(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var c commander.Commander

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, c, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h CommanderHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	id := chi.URLParam(r, "allianceid")

	cs, err := commander.ListByAlliance(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, cs, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h CommanderHandler) Merge(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Merge Commanders not implemented ...\n"))
}

func (h CommanderHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var err error
	var c commander.Commander

	name := chi.URLParam(r, "name")
	err = c.GetByName(name)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, c, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h CommanderHandler) Update(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var c commander.Commander

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, c, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h CommanderHandler) AddData(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var c commander.Commander
	c.Id = chi.URLParam(r, "id")

	date := chi.URLParam(r, "date")

	err := c.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var d commander.Data
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.AddData(date, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, c, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CommanderRoutes() chi.Router {
	r := chi.NewRouter()

	h := CommanderHandler{}
	r.Get("/", h.List)
	r.Get("/a/{allianceid}", h.ListMembers)
	r.Get("/n/{name}", h.GetByName)
	r.Post("/{server}", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Put("/{id}/data/{date}", h.AddData)
	r.Put("/{id}/m/{id2}", h.Merge)

	return r
}
