package api

import (
	"github.com/gorilla/mux"
	"net/http"
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

type Context struct {
	TremorRepo
	MedicineRepo
	ExerciseRepo
	Reboot chan struct{}
}

func NewRouter(ctx *Context) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/tremors", getTremorsSince(ctx.TremorRepo)).
		Queries("since", "{since}").Methods(http.MethodGet)
	router.HandleFunc("/api/tremors", getTremors(ctx.TremorRepo)).Methods(http.MethodGet)
	router.HandleFunc("/api/tremors", addTremor(ctx.TremorRepo)).Methods(http.MethodPost)
	router.HandleFunc("/api/meds/{mid}", getMedicine(ctx.MedicineRepo)).Methods(http.MethodGet)
	router.HandleFunc("/api/meds", getMedicinesForDate(ctx.MedicineRepo)).
		Queries("date", "{date}").Methods(http.MethodGet)
	router.HandleFunc("/api/meds", updateMedicine(ctx.MedicineRepo)).Methods(http.MethodPut)
	router.HandleFunc("/api/meds", getMedicines(ctx.MedicineRepo)).Methods(http.MethodGet)
	router.HandleFunc("/api/meds", addMedicine(ctx.MedicineRepo)).Methods(http.MethodPost)
	router.HandleFunc("/api/exercises/{eid}", getExercise(ctx.ExerciseRepo)).Methods(http.MethodGet)
	router.HandleFunc("/api/exercises", updateExercise(ctx.ExerciseRepo)).Methods(http.MethodPut)
	router.HandleFunc("/api/exercises", getExercises(ctx.ExerciseRepo)).Methods(http.MethodGet)
	router.HandleFunc("/api/exercises", addExercise(ctx.ExerciseRepo)).Methods(http.MethodPost)
	router.HandleFunc("/api/update", update(ctx.Reboot)).Methods(http.MethodPost)
	return router
}
