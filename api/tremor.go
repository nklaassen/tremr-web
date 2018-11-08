package api

import (
	"encoding/json"
	"net/http"
	"github.com/nklaassen/tremr-web/datastore"
	"log"
)

func getTremors(tremorRepo datastore.TremorRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, h *http.Request) {
		tremors, err := tremorRepo.GetAll()
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tremors)
	}
}

func addTremor(tremorRepo datastore.TremorRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tremor datastore.Tremor
		if err := json.NewDecoder(r.Body).Decode(&tremor); err != nil {
			log.Print("decode error: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := tremorRepo.Add(&tremor); err != nil {
			log.Print("datastore error: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}