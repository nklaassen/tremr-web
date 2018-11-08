package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nklaassen/tremr-web/api"
	"github.com/nklaassen/tremr-web/datastore/sql"
	"log"
	"net/http"
	"os"
)

func main() {
	apiContext := &api.Context{sql.CreateTremorRepo()}

	apiserver := api.NewRouter(apiContext)

	fileserver := http.FileServer(http.Dir("www"))

	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(apiserver).Methods("GET", "POST")
	router.PathPrefix("/").Handler(fileserver).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
