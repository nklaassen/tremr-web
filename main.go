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
	"os/signal"
	"syscall"
)

func createApiContext(db *sqlx.DB) *api.Context {
	tremorRepo, err := database.NewTremorRepo(db)
	if err != nil {
		log.Fatal("Failed to create TremorRepo: ", err)
	}

	medicineRepo, err := database.NewMedicineRepo(db)
	if err != nil {
		log.Fatal("Failed to create MedicineRepo: ", err)
	}

	exerciseRepo, err := database.NewExerciseRepo(db)
	if err != nil {
		log.Fatal("Failed to create ExerciseRepo: ", err)
	}

	return &api.Context{
		TremorRepo:   tremorRepo,
		MedicineRepo: medicineRepo,
		ExerciseRepo: exerciseRepo}
}

func main() {
	c := make(chan os.Signal)

	// Open database
	db, err := sqlx.Open("sqlite3", "db.sqlite3?_journal_mode=WAL")
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	// Close the database and exit on SIGINT and SIGTERM
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Print("Caught signal, closing database and exiting")
		db.Close()
		os.Exit(1)
	}()

	// Create API server
	apiContext := createApiContext(db)
	apiserver := api.NewRouter(apiContext)

	// Create fileserver out of www/ directory
	fileserver := http.FileServer(http.Dir("www"))

	// Set up router
	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(apiserver).Methods("GET", "POST")
	router.PathPrefix("/").Handler(fileserver).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// Start the server on port 8080, this should never return
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
