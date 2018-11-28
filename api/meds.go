package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type Schedule struct {
	Mo bool `json:"mo"`
	Tu bool `json:"tu"`
	We bool `json:"we"`
	Th bool `json:"th"`
	Fr bool `json:"fr"`
	Sa bool `json:"sa"`
	Su bool `json:"su"`
}

type Medicine struct {
	MID       int64  `json:"mid"`
	UID       int64  `json:"uid"`
	Name      string `json:"name"`
	Dosage    string `json:"dosage"`
	Schedule  `json:"schedule"`
	Reminder  bool       `json:"reminder"`
	StartDate time.Time  `json:"startdate"`
	EndDate   *time.Time `json:"enddate"`
}

type MedicineRepo interface {
	Add(uid int64, med *Medicine) error
	GetAll(uid int64) ([]Medicine, error)
	Get(uid, mid int64) (Medicine, error)
	GetForDate(uid int64, date time.Time) ([]Medicine, error)
	Update(uid int64, med *Medicine) error
}

func medsRouter(repo MedicineRepo) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/meds/{mid}", updateMedicine(repo)).Methods(http.MethodPut)
	router.Handle("/meds/{mid}", getMedicine(repo)).Methods(http.MethodGet)
	router.Handle("/meds", getMedicinesForDate(repo)).Queries("date", "{date}").Methods(http.MethodGet)
	router.Handle("/meds", getMedicines(repo)).Methods(http.MethodGet)
	router.Handle("/meds", addMedicine(repo)).Methods(http.MethodPost)
	return router
}

func getMedicines(medicineRepo MedicineRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		medicines, err := medicineRepo.GetAll(uid)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(medicines)
		return nil
	}
}

func addMedicine(medicineRepo MedicineRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		// decode medicine from json in body of request
		var medicine Medicine
		if err := json.NewDecoder(r.Body).Decode(&medicine); err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}
		if medicine.Name == "" || medicine.Dosage == "" || medicine.Schedule == (Schedule{}) {
			return HandlerError{errors.New("must populate name, dosage, schedule"), http.StatusBadRequest}
		}
		return medicineRepo.Add(uid, &medicine)
	}
}

func getMedicine(medicineRepo MedicineRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		// get mid from url
		vars := mux.Vars(r)
		mid, err := strconv.ParseInt(vars["mid"], 10, 64)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		medicine, err := medicineRepo.Get(uid, mid)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(medicine)
		return nil
	}
}

func getMedicinesForDate(medicineRepo MedicineRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		// get date from url
		timestring := r.FormValue("date")
		timestamp, err := time.Parse(time.RFC3339, timestring)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		medicines, err := medicineRepo.GetForDate(uid, timestamp)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(medicines)
		return nil
	}
}

func updateMedicine(medicineRepo MedicineRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		uid := r.Context().Value("uid").(int64)
		// get mid from url
		vars := mux.Vars(r)
		mid, err := strconv.ParseInt(vars["mid"], 10, 64)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}
		// parse medicine from json in body of request
		var medicine Medicine
		if err := json.NewDecoder(r.Body).Decode(&medicine); err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}
		if mid != medicine.MID {
			return HandlerError{errors.New("mid in url and body do not match"), http.StatusBadRequest}
		}
		if medicine.MID == 0 || medicine.Name == "" || medicine.Dosage == "" ||
			medicine.Schedule == (Schedule{}) {
			return HandlerError{errors.New("must populate all fields for update"), http.StatusBadRequest}
		}
		return medicineRepo.Update(uid, &medicine)
	}
}
