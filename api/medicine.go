package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Medicine struct {
	MID       *int64  `json:"mid"`
	Name      *string `json:"name"`
	Dosage    *string `json:"dosage"`
	*Schedule `json:"schedule"`
	Reminder  *bool      `json:"reminder"`
	StartDate *time.Time `json:"startdate"`
	EndDate   *time.Time `json:"enddate"`
}

type MedicineRepo interface {
	Add(*Medicine) error
	GetAll() ([]Medicine, error)
	Get(int64) (Medicine, error)
	Update(*Medicine) error
}

func getMedicines(medicineRepo MedicineRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, h *http.Request) {
		medicines, err := medicineRepo.GetAll()
		if err != nil {
			log.Print(err)
			http.Error(w, "failed to get medicines from database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(medicines)
	}
}

func addMedicine(medicineRepo MedicineRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var medicine Medicine
		if err := json.NewDecoder(r.Body).Decode(&medicine); err != nil {
			log.Print("decode error: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if medicine.Name == nil || medicine.Dosage == nil || medicine.Schedule == nil {
			log.Print("invalid json request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := medicineRepo.Add(&medicine); err != nil {
			log.Print("database error: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func getMedicine(medicineRepo MedicineRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mid, err := strconv.ParseInt(vars["mid"], 10, 64)
		if err != nil {
			http.Error(w, "invalid medicine id in url", http.StatusBadRequest)
			return
		}

		medicine, err := medicineRepo.Get(mid)
		if err != nil {
			log.Print("database error:", err)
			http.Error(w, "failed to get medicine from database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(medicine)
	}
}

func updateMedicine(medicineRepo MedicineRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var medicine Medicine
		if err := json.NewDecoder(r.Body).Decode(&medicine); err != nil {
			log.Print("decode error: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if medicine.MID == nil || medicine.Name == nil || medicine.Dosage == nil ||
			medicine.Schedule == nil || medicine.Reminder == nil {
			log.Print("invalid json request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := medicineRepo.Update(&medicine); err != nil {
			log.Print("database error:", err)
			http.Error(w, "error updating database", http.StatusInternalServerError)
			return
		}
	}
}
