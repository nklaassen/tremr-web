package main

import (
	"context"
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

func createApiContext(db *sqlx.DB, reboot chan struct{}) *api.Context {
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
		ExerciseRepo: exerciseRepo,
		Reboot:       reboot,
	}
}

func serve(portNum string, reboot chan struct{}, shutdown chan struct{}) {
	// Open database
	db, err := sqlx.Open("sqlite3", "db.sqlite3?_journal_mode=WAL")
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	// Create API server
	apiContext := createApiContext(db, reboot)
	apiserver := api.NewRouter(apiContext)

	// Create fileserver out of www/ directory
	fileserver := http.FileServer(http.Dir("www"))

	// Set up router
	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(apiserver).Methods("GET", "POST", "PUT")
	router.PathPrefix("/").Handler(fileserver).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// Create server
	srv := &http.Server{
		Addr:    ":" + portNum,
		Handler: loggedRouter,
	}

	// Start the webserver in a goroutine so we can wait on the shutdown channel
	go func() {
		log.Printf("starting tremr-web on port %v\n", portNum)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Error starting server: %v", err)
			close(shutdown)
		}
	}()
	<-shutdown
	srv.Shutdown(context.Background())
}

func main() {
	// Get port num from cmd line arg, default to 8080
	portNum := "8080"
	if len(os.Args) > 1 {
		portNum = os.Args[1]
	}

	// if this channel is closed, restart the go binary
	reboot := make(chan struct{})
	// close this channel to shutdown the web server
	shutdown := make(chan struct{})

	// Start a goroutine which shuts down the webserver on a signal,
	//	or a send to the reboot channel
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		select {
		case <-c:
			log.Print("Caught signal, shutting down server")
		case <-reboot:
			log.Print("Shutting down server and restarting go binary")
		}
		close(shutdown)
	}()

	// this blocks until the webserver shuts down
	serve(portNum, reboot, shutdown)

	// if the reboot channel was closed, execve the (probably updated) go binary (see ./api/update)
	select {
	case <-reboot:
		if err := syscall.Exec(os.Args[0], os.Args, os.Environ()); err != nil {
			log.Print("Failed to restart server", err)
		}
	default:
	}
}
