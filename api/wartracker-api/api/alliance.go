package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"wartracker/pkg/alliance"

	"github.com/go-chi/chi/v5"
)

type AllianceHandler struct {
}

func (h AllianceHandler) ListAlliances(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")
	data := GetQueryBool(r, "data")

	as, err := alliance.List(data)
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

func (h AllianceHandler) GetAlliance(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")
	data := GetQueryBool(r, "data")
	latest := GetQueryBool(r, "latest")

	var err error
	var a alliance.Alliance
	a.Id = chi.URLParam(r, "id")
	if latest {
		err = a.Get(latest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if data {
		err = a.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = a.GetData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, ErrAliianceNotFound.Error(), http.StatusNotFound)
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

func (h AllianceHandler) CreateAlliance(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var a alliance.Alliance

	s := chi.URLParam(r, "server")
	si, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.Server = si

	err = json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.Create()
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

func (h AllianceHandler) UpdateAlliance(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update alliance not implemented ...\n"))
}

func (h AllianceHandler) AddAllianceData(w http.ResponseWriter, r *http.Request) {
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
	a.DataMap[d.Date] = d
	err = a.AddData(d.Date)
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

	allianceHandler := AllianceHandler{}
	r.Get("/", allianceHandler.ListAlliances)
	r.Post("/{server}", allianceHandler.CreateAlliance)
	r.Get("/{id}", allianceHandler.GetAlliance)
	r.Put("/{id}", allianceHandler.UpdateAlliance)
	r.Put("/{id}/data", allianceHandler.AddAllianceData)

	return r
}
