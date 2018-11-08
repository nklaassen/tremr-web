package api

import (
	"github.com/nklaassen/tremr-web/datastore"
	"github.com/gorilla/mux"
)

type Context struct {
	TremorRepo datastore.TremorRepo
}

func NewRouter(ctx *Context) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/tremors", getTremors(ctx.TremorRepo)).Methods("GET")
	router.HandleFunc("/api/tremors", addTremor(ctx.TremorRepo)).Methods("POST")
	return router
}
