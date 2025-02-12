package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"wartracker/pkg/commander"

	"github.com/go-chi/chi/v5"
)

type CommanderHandler struct {
}

func (h CommanderHandler) ListCommanders(w http.ResponseWriter, r *http.Request) {
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

func (h CommanderHandler) GetCommander(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")
	latest := GetQueryBool(r, "latest")

	var err error
	var c commander.Commander
	c.Id = chi.URLParam(r, "id")
	err = c.GetById(c.Id, latest)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, ErrAliianceNotFound.Error(), http.StatusNotFound)
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

func (h CommanderHandler) CreateCommander(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var c commander.Commander

	s := chi.URLParam(r, "server")
	si, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.Server = si

	err = json.NewDecoder(r.Body).Decode(&c)
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

func (h CommanderHandler) ListAllianceCommanders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List alliance Commanders not implemented ...\n"))
}

func (h CommanderHandler) MergeCommanders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List alliance Commanders not implemented ...\n"))
}

func (h CommanderHandler) GetCommanderByName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Commander by name not implemented ...\n"))
}

func (h CommanderHandler) UpdateCommander(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update Commander not implemented ...\n"))
}

func (h CommanderHandler) AddCommanderData(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var c commander.Commander
	c.Id = chi.URLParam(r, "id")
	err := c.GetById(c.Id)
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
	err = c.AddData()
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

	CommanderHandler := CommanderHandler{}
	r.Get("/", CommanderHandler.ListCommanders)
	r.Get("/a/{allianceid}", CommanderHandler.ListAllianceCommanders)
	r.Get("/n/{name}", CommanderHandler.GetCommanderByName)
	r.Post("/{server}", CommanderHandler.CreateCommander)
	r.Get("/{id}", CommanderHandler.GetCommander)
	r.Put("/{id}", CommanderHandler.UpdateCommander)
	r.Put("/{id}/data", CommanderHandler.AddCommanderData)
	r.Put("/{id}/m/{id2}", CommanderHandler.MergeCommanders)

	return r
}
