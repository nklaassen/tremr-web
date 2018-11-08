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

type TremorRepo interface {
	Add(*Tremor) error
	GetAll() ([]Tremor, error)
}

type Context struct {
	TremorRepo TremorRepo
}

func NewRouter(ctx *Context) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/tremors", getTremors(ctx.TremorRepo)).Methods("GET")
	router.HandleFunc("/api/tremors", addTremor(ctx.TremorRepo)).Methods("POST")
	return router
}
