package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nklaassen/tremr-web/api"
	"github.com/nklaassen/tremr-web/database"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := sqlx.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	tremorRepo, err := database.NewTremorRepo(db)
	if err != nil {
		log.Fatal("Failed to create TremorRepo: ", err)
	}

	medicineRepo, err := database.NewMedicineRepo(db)
	if err != nil {
		log.Fatal("Failed to create MedicineRepo: ", err)
	}

	apiContext := &api.Context{tremorRepo, medicineRepo}
	apiserver := api.NewRouter(apiContext)

	fileserver := http.FileServer(http.Dir("www"))

	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(apiserver).Methods("GET", "POST")
	router.PathPrefix("/").Handler(fileserver).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
