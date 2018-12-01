package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type Exercise struct {
	EID       int64  `json:"eid"`
	UID       int64  `json:"uid"`
	Name      string `json:"name"`
	Unit      string `json:"unit"`
	Schedule  `json:"schedule"`
	Reminder  bool       `json:"reminder"`
	StartDate time.Time  `json:"startdate"`
	EndDate   *time.Time `json:"enddate"`
}

type ExerciseRepo interface {
	Add(uid int64, exer *Exercise) error
	GetAll(uid int64) ([]Exercise, error)
	Get(uid, eid int64) (Exercise, error)
	GetForDate(uid int64, date time.Time) ([]Exercise, error)
	Update(uid int64, exer *Exercise) error
}

func exercisesRouter(repo ExerciseRepo) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/exercises/{eid}", updateExercise(repo)).Methods(http.MethodPut)
	r.Handle("/exercises/{eid}", getExercise(repo)).Methods(http.MethodGet)
	r.Handle("/exercises", getExercisesForDate(repo)).Queries("date", "{date}").Methods(http.MethodGet)
	r.Handle("/exercises", getExercises(repo)).Queries("uid", "{uid}").Methods(http.MethodGet)
	r.Handle("/exercises", getExercises(repo)).Methods(http.MethodGet)
	r.Handle("/exercises", addExercise(repo)).Methods(http.MethodPost)
	return r
}

func getExercises(exerciseRepo ExerciseRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var forUid int64
		var err error
		// get uid from token, added to context by authMiddleware
		tokenUid := r.Context().Value("uid").(int64)

		// if the logged in user is trying to read the exercises of another user
		requestedUidString := r.FormValue("uid")
		if requestedUidString == "" {
			// default to the logged in user
			forUid = tokenUid
		} else {
			// else get the exercises for the uid requested in the url
			// FIXME: need to add authentication here - make sure that user has actually shared with us
			forUid, err = strconv.ParseInt(requestedUidString, 10, 64)
			if err != nil {
				return HandlerError{err, http.StatusBadRequest}
			}
		}

		exercises, err := exerciseRepo.GetAll(forUid)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
		return nil
	}
}

func addExercise(exerciseRepo ExerciseRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		// decode exercise from json in body of request
		var exercise Exercise
		if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}
		if exercise.Name == "" || exercise.Unit == "" || exercise.Schedule == (Schedule{}) {
			return HandlerError{errors.New("must populate name, unit, schedule"), http.StatusBadRequest}
		}
		return exerciseRepo.Add(uid, &exercise)
	}
}

func getExercise(exerciseRepo ExerciseRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		// get eid from url
		vars := mux.Vars(r)
		eid, err := strconv.ParseInt(vars["eid"], 10, 64)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		exercise, err := exerciseRepo.Get(uid, eid)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercise)
		return nil
	}
}

func getExercisesForDate(exerciseRepo ExerciseRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		// get date from url
		timestring := r.FormValue("date")
		timestamp, err := time.Parse(time.RFC3339, timestring)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		exercises, err := exerciseRepo.GetForDate(uid, timestamp)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
		return nil
	}
}

func updateExercise(exerciseRepo ExerciseRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		// get eid from url
		vars := mux.Vars(r)
		eid, err := strconv.ParseInt(vars["eid"], 10, 64)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}
		// parse exercise from json in body of request
		var exercise Exercise
		if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
			return err
		}
		if eid != exercise.EID {
			return HandlerError{errors.New("eid in url and body do not match"), http.StatusBadRequest}
		}
		if exercise.EID == 0 || exercise.Name == "" || exercise.Unit == "" ||
			exercise.Schedule == (Schedule{}) {
			return HandlerError{errors.New("must populate all fields for update"), http.StatusBadRequest}
		}
		return exerciseRepo.Update(uid, &exercise)
	}
}
