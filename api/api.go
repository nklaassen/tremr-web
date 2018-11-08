package api

import (
	"github.com/gorilla/mux"
)

type Tremor struct {
	Tid       int `json:"tid"`
	Resting   int `json:"resting"`
	Postural  int `json:"postural"`
	Completed bool `json:"completed"`
	Date      string `json:"date"`
}

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
	MID int `json:"mid"`
	Name string `json:"name"`
	Dosage string `json:"dosage"`
	Schedule `json:"schedule"`
	Reminder bool `json:"reminder"`
	StartDate string `json:"startdate"`
	EndDate *string `json:"enddate"`
}

type TremorRepo interface {
	Add(*Tremor) error
	GetAll() ([]Tremor, error)
}

type MedicineRepo interface {
	Add(*Medicine) error
	GetAll() ([]Medicine, error)
}

type Context struct {
	TremorRepo
	MedicineRepo
}

func NewRouter(ctx *Context) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/tremors", getTremors(ctx.TremorRepo)).Methods("GET")
	router.HandleFunc("/api/tremors", addTremor(ctx.TremorRepo)).Methods("POST")
	router.HandleFunc("/api/meds", getMedicines(ctx.MedicineRepo)).Methods("GET")
	router.HandleFunc("/api/meds", addMedicine(ctx.MedicineRepo)).Methods("POST")
	return router
}
