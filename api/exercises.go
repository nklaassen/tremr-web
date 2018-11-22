package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Exercise struct {
	EID       *int64  `json:"eid"`
	Name      *string `json:"name"`
	Unit      *string `json:"unit"`
	*Schedule `json:"schedule"`
	Reminder  *bool      `json:"reminder"`
	StartDate *time.Time `json:"startdate"`
	EndDate   *time.Time `json:"enddate"`
}

type ExerciseRepo interface {
	Add(*Exercise) error
	GetAll() ([]Exercise, error)
	Get(int64) (Exercise, error)
	GetForDate(time.Time) ([]Exercise, error)
	Update(*Exercise) error
}

func getExercises(exerciseRepo ExerciseRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, h *http.Request) {
		exercises, err := exerciseRepo.GetAll()
		if err != nil {
			log.Print(err)
			http.Error(w, "failed to get exercises from database", http.StatusInternalServerError)
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

func getExercise(exerciseRepo ExerciseRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		eid, err := strconv.ParseInt(vars["eid"], 10, 64)
		if err != nil {
			http.Error(w, "invalid exercise id in url", http.StatusBadRequest)
			return
		}

		exercise, err := exerciseRepo.Get(eid)
		if err != nil {
			log.Print("database error:", err)
			http.Error(w, "failed to get exercise from database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercise)
	}
}

func getExercisesForDate(exerciseRepo ExerciseRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		timestring := r.FormValue("date")
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

		exercises, err := exerciseRepo.GetForDate(timestamp)
		if err != nil {
			log.Print("database error:", err)
			http.Error(w, "failed to get exercises from database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
	}
}

func updateExercise(exerciseRepo ExerciseRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var exercise Exercise
		if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
			log.Print("decode error: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if exercise.EID == nil || exercise.Name == nil || exercise.Unit == nil ||
			exercise.Schedule == nil || exercise.Reminder == nil {
			log.Print("invalid json request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := exerciseRepo.Update(&exercise); err != nil {
			log.Print("database error:", err)
			http.Error(w, "error updating database", http.StatusInternalServerError)
			return
		}
	}
}
