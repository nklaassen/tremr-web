package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"strconv"
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
	router.Handle("/tremors", getTremors(repo)).Queries("uid", "{uid}").Methods(http.MethodGet)
	router.Handle("/tremors", getTremors(repo)).Methods(http.MethodGet)
	router.Handle("/tremors", addTremor(repo)).Methods(http.MethodPost)
	return router
}

func getTremors(tremorRepo TremorRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var forUid int64
		var err error
		// get uid from token, added to context by authMiddleware
		tokenUid := r.Context().Value("uid").(int64)

		// if the logged in user is trying to read the tremors of another user
		requestedUidString := r.FormValue("uid")
		if requestedUidString == "" {
			// default to the logged in user
			forUid = tokenUid
		} else {
			// else get the tremors for the uid requested in the url
			// FIXME: need to add authentication here - make sure that user has actually shared with us
			forUid, err = strconv.ParseInt(requestedUidString, 10, 64)
			if err != nil {
				return HandlerError{err, http.StatusBadRequest}
			}
		}

		tremors, err := tremorRepo.GetAll(forUid)
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
		// get uid from token, added to context by authMiddleware
		tokenUid := r.Context().Value("uid").(int64)

		// get since timestamp from url
		timestring := r.FormValue("since")
		timestamp, err := time.Parse(time.RFC3339, timestring)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		// get all tremors for the logged in user since the requested date
		tremors, err := tremorRepo.GetSince(tokenUid, timestamp)
		if err != nil {
			return HandlerError{err, http.StatusInternalServerError}
		}

		// return the tremors in a JSON array
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tremors)
		return nil
	}
}

func addTremor(tremorRepo TremorRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		tokenUid := r.Context().Value("uid").(int64)

		// decode the tremor from JSON in the request body
		var tremor Tremor
		if err := json.NewDecoder(r.Body).Decode(&tremor); err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		// add the tremor to the db and return any error
		return tremorRepo.Add(tokenUid, &tremor)
	}
}
