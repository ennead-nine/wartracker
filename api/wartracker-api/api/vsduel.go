package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"wartracker/pkg/vsduel"

	"github.com/go-chi/chi/v5"
)

func init() {
	initVsDuel()
}

func initVsDuel() {
	var err error
	vsduel.Days, err = vsduel.GetDays()
	if err != nil {
		panic(err)
	}
}

type VsDuelHandler struct {
}

func (h VsDuelHandler) ListVsDuel(w http.ResponseWriter, r *http.Request) {
	http.Error(w, ErrNotImplemented.Error(), http.StatusNotImplemented)
}

func (h VsDuelHandler) Get(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	vs, err := vsduel.Get(chi.URLParam(r, "id"))
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

	vd, err := vsduel.Get(chi.URLParam(r, "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = vd.GetWeeks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	week, err := strconv.Atoi(chi.URLParam(r, "week"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if k, ok := vd.Weeks[week]; ok {
		err = k.StartDay(chi.URLParam(r, "day"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = k.ScanPointsRanking(z, chi.URLParam(r, "day"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		vd.Weeks[week] = k
	}

	jsonOutput(w, vd.Weeks[week].Data[chi.URLParam(r, "day")], indent)
}

func (h VsDuelHandler) StartVsWeek(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var k vsduel.Week
	err := json.NewDecoder(r.Body).Decode(&k)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vs, err := vsduel.Get(k.VsDuelId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = vs.StartWeek(k.WeekNumber, k.AllianceIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonOutput(w, vs.Weeks[k.WeekNumber], indent)
}

func VsDuelRoutes() chi.Router {
	r := chi.NewRouter()

	vsDuelHandler := VsDuelHandler{}
	r.Get("/", vsDuelHandler.ListVsDuel)
	r.Post("/", vsDuelHandler.CreateVsDuel)
	r.Post("/week", vsDuelHandler.StartVsWeek)
	r.Post("/scan/{id}/{week}/{day}", vsDuelHandler.ScanVsDay)
	r.Get("/{id}", vsDuelHandler.Get)
	r.Put("/{id}", vsDuelHandler.UpdateVsDuel)

	return r
}
