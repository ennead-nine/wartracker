package api

import (
	"encoding/json"
	"io"
	"net/http"
	"wartracker/pkg/donation"

	"github.com/go-chi/chi/v5"
)

type DonationHandler struct {
}

func (h DonationHandler) Create(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var d donation.Donation

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = donation.AddDonation(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = jsonOutput(w, d, indent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h DonationHandler) Scan(w http.ResponseWriter, r *http.Request) {
	indent := GetQueryBool(r, "indent")

	var err error
	var z []byte

	z, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	dm := make(donation.DonationMap)
	badimg, err := dm.ScanDonations(z, chi.URLParam(r, "date"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonOutput(w, dm, indent)
	jsonOutput(w, badimg, indent)
}

func DonationRoutes() chi.Router {
	r := chi.NewRouter()

	h := DonationHandler{}
	r.Post("/", h.Create)
	r.Post("/scan/{date}", h.Scan)

	return r
}
