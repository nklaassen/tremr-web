package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Tremor struct {
	TID      int64     `json:"tid"`
	UID      int64     `json:"uid"`
	Resting  int       `json:"resting"`
	Postural int       `json:"postural"`
	Date     time.Time `json:"date"`
}

type TremorRepo interface {
	Add(uid int64, tremor *Tremor) error
	GetAll(uid int64) ([]Tremor, error)
	GetSince(uid int64, since time.Time) ([]Tremor, error)
}

func tremorsRouter(repo TremorRepo) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/tremors", getTremorsSince(repo)).Queries("since", "{since}").Methods(http.MethodGet)
	router.Handle("/tremors", getTremors(repo)).Methods(http.MethodGet)
	router.Handle("/tremors", addTremor(repo)).Methods(http.MethodPost)
	return router
}

func getTremors(tremorRepo TremorRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.Context().Value("uid").(int64)
		tremors, err := tremorRepo.GetAll(uid)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tremors)
		return nil
	}
}

func getTremorsSince(tremorRepo TremorRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.Context().Value("uid").(int64)
		timestring := r.FormValue("since")
		timestamp, err := time.Parse(time.RFC3339, timestring)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}
		tremors, err := tremorRepo.GetSince(uid, timestamp)
		if err != nil {
			return HandlerError{err, http.StatusInternalServerError}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tremors)
		return nil
	}
}

func addTremor(tremorRepo TremorRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.Context().Value("uid").(int64)
		var tremor Tremor
		if err := json.NewDecoder(r.Body).Decode(&tremor); err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}
		return tremorRepo.Add(uid, &tremor)
	}
}
