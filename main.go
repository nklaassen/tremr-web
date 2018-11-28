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

func serve(portNum string, reboot chan struct{}, shutdown chan struct{}) {
	// Open raw database
	db, err := sqlx.Open("sqlite3", "db.sqlite3?_journal_mode=WAL")
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	// Get datastore
	ds, err := database.GetDataStore(db)
	if err != nil {
		log.Fatal(err)
	}

	// Create API server
	apiserver := api.NewRouter(&api.Env{ds, reboot})

	// Create fileserver out of www/ directory
	fileserver := http.FileServer(http.Dir("www"))

	// Set up router
	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(http.StripPrefix("/api", apiserver))
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		select {
		case <-sig:
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
