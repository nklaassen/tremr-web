package api

import (
	"github.com/gorilla/mux"
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
}

func NewRouter(ctx *Context) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/tremors", getTremorsSince(ctx.TremorRepo)).
		Queries("since", "{since}").Methods("GET")
	router.HandleFunc("/api/tremors", getTremors(ctx.TremorRepo)).Methods("GET")
	router.HandleFunc("/api/tremors", addTremor(ctx.TremorRepo)).Methods("POST")
	router.HandleFunc("/api/meds", getMedicines(ctx.MedicineRepo)).Methods("GET")
	router.HandleFunc("/api/meds", addMedicine(ctx.MedicineRepo)).Methods("POST")
	router.HandleFunc("/api/exercises", getExercises(ctx.ExerciseRepo)).Methods("GET")
	router.HandleFunc("/api/exercises", addExercise(ctx.ExerciseRepo)).Methods("POST")
	return router
}
