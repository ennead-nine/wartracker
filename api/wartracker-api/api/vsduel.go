package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	//	err = vsduel.InitDays()
	//	if err != nil {
	//		panic(err)
	//	}
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

	var badimg []string
	if k, ok := vd.Weeks[week]; ok {
		err = k.StartDay(chi.URLParam(r, "day"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		badimg, err = k.ScanPointsRanking(z, chi.URLParam(r, "day"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		vd.Weeks[week] = k
	}

	jsonOutput(w, vd.Weeks[week].Data[chi.URLParam(r, "day")], indent)
	jsonOutput(w, badimg, indent)
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

func (h VsDuelHandler) CalculateAllianceData(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	vd, err := vsduel.Get(chi.URLParam(r, "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			lerr := fmt.Errorf("vsduel %s not found: %w", chi.URLParam(r, "id"), err)
			http.Error(w, lerr.Error(), http.StatusNotFound)
			return
		} else {
			lerr := fmt.Errorf("failed to get vsduel %s: %w", chi.URLParam(r, "id"), err)
			http.Error(w, lerr.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = vd.GetWeeks()
	if err != nil {
		lerr := fmt.Errorf("failed to get weeks for %s: %w", chi.URLParam(r, "id"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}

	week, err := strconv.Atoi(chi.URLParam(r, "week"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	k := vd.Weeks[week]
	err = k.GetDays()
	if err != nil {
		lerr := fmt.Errorf("failed to get days for week %s: %w", chi.URLParam(r, "week"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}
	d := k.Data[chi.URLParam(r, "day")]
	err = d.GetCommanderData()
	if err != nil {
		lerr := fmt.Errorf("failed to get commander data for %s: %w", chi.URLParam(r, "day"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}
	err = d.CalculateAllianceData()
	if err != nil {
		lerr := fmt.Errorf("failed to calculate alliance data for %s: %w", chi.URLParam(r, "day"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}

	jsonOutput(w, d.AllianceData, indent)
}

func (h VsDuelHandler) UpdateRanks(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	vd, err := vsduel.Get(chi.URLParam(r, "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			lerr := fmt.Errorf("vsduel %s not found: %w", chi.URLParam(r, "id"), err)
			http.Error(w, lerr.Error(), http.StatusNotFound)
			return
		} else {
			lerr := fmt.Errorf("failed to get vsduel %s: %w", chi.URLParam(r, "id"), err)
			http.Error(w, lerr.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = vd.GetWeeks()
	if err != nil {
		lerr := fmt.Errorf("failed to get weeks for %s: %w", chi.URLParam(r, "id"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}

	week, err := strconv.Atoi(chi.URLParam(r, "week"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	k := vd.Weeks[week]
	err = k.GetDays()
	if err != nil {
		lerr := fmt.Errorf("failed to get days for week %s: %w", chi.URLParam(r, "week"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}

	d := k.Data[chi.URLParam(r, "day")]
	err = d.GetCommanderData()
	if err != nil {
		lerr := fmt.Errorf("failed to get commander data for %s: %w", chi.URLParam(r, "day"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}

	cd := d.CommanderData
	err = cd.UpdateRanks()
	if err != nil {
		lerr := fmt.Errorf("failed to update rankings  data for %s: %w", chi.URLParam(r, "day"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}

	jsonOutput(w, d.CommanderData, indent)
}

func (h VsDuelHandler) MergeCommanderData(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	vd, err := vsduel.Get(chi.URLParam(r, "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			lerr := fmt.Errorf("vsduel %s not found: %w", chi.URLParam(r, "id"), err)
			http.Error(w, lerr.Error(), http.StatusNotFound)
			return
		} else {
			lerr := fmt.Errorf("failed to get vsduel %s: %w", chi.URLParam(r, "id"), err)
			http.Error(w, lerr.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = vd.GetWeeks()
	if err != nil {
		lerr := fmt.Errorf("failed to get weeks for %s: %w", chi.URLParam(r, "id"), err)
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
		return
	}

	week, err := strconv.Atoi(chi.URLParam(r, "week"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	k := vd.Weeks[week]

	err = k.MergeCommanderData(chi.URLParam(r, "src"), chi.URLParam(r, "dst"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonOutput(w, k.Data, indent)
}

func VsDuelRoutes() chi.Router {
	r := chi.NewRouter()

	h := VsDuelHandler{}
	r.Get("/", h.ListVsDuel)
	r.Post("/", h.CreateVsDuel)
	r.Post("/week", h.StartVsWeek)
	r.Post("/{id}/{week}/{day}/scan", h.ScanVsDay)
	r.Put("/{id}/{week}/{day}/calc", h.CalculateAllianceData)
	r.Put("/{id}/{week}/{day}/ranks", h.UpdateRanks)
	r.Put("/merge/{id}/{week}/{src}/{dst}", h.MergeCommanderData)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.UpdateVsDuel)

	return r
}
