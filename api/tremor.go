package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Tremor struct {
	Tid      int       `json:"tid"`
	Resting  *int      `json:"resting"`
	Postural *int      `json:"postural"`
	Date     *time.Time `json:"date"`
}

type TremorRepo interface {
	Add(*Tremor) error
	GetAll() ([]Tremor, error)
	GetSince(time.Time) ([]Tremor, error)
}

func getTremors(tremorRepo TremorRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

func getTremorsSince(tremorRepo TremorRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		timestring := r.FormValue("since")
		if timestring == "" {
			log.Print("invalid query")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		timestamp, err := time.Parse(time.RFC3339, timestring)
		if err != nil {
			log.Print("invalid timestamp")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print("timestamp: ", timestamp)
		tremors, err := tremorRepo.GetSince(timestamp)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tremors)
	}
}

func addTremor(tremorRepo TremorRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tremor Tremor
		if err := json.NewDecoder(r.Body).Decode(&tremor); err != nil {
			log.Print("decode error: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if tremor.Resting == nil || tremor.Postural == nil {
			log.Print("invalid json request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := tremorRepo.Add(&tremor); err != nil {
			log.Print("database error: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
