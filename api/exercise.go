package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Exercise struct {
	EID       int     `json:"eid"`
	Name      *string `json:"name"`
	Unit      *string `json:"unit"`
	*Schedule `json:"schedule"`
	Reminder  bool       `json:"reminder"`
	StartDate *time.Time `json:"startdate"`
	EndDate   *time.Time `json:"enddate"`
}

type ExerciseRepo interface {
	Add(*Exercise) error
	GetAll() ([]Exercise, error)
}

func getExercises(exerciseRepo ExerciseRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, h *http.Request) {
		exercises, err := exerciseRepo.GetAll()
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
	}
}

func addExercise(exerciseRepo ExerciseRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var exercise Exercise
		if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
			log.Print("decode error: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if exercise.Name == nil || exercise.Unit == nil || exercise.Schedule == nil {
			log.Print("invalid json request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := exerciseRepo.Add(&exercise); err != nil {
			log.Print("database error: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
