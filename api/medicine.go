package api

import (
	"encoding/json"
	"net/http"
	"log"
)

func getMedicines(medicineRepo MedicineRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, h *http.Request) {
		medicines, err := medicineRepo.GetAll()
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
