package main

import (
	"fmt"
	//	"github.com/gorilla/handlers"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nklaassen/tremr-web/api"
	"github.com/nklaassen/tremr-web/database"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var apiContext *api.Context
var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Open("sqlite3", "test_db.sqlite3?_journal_mode=WAL")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`drop table if exists tremors;
		drop table if exists medicines;
		drop table if exists exercises`)
	if err != nil {
		panic(err)
	}

	tremorRepo, err := database.NewTremorRepo(db)
	if err != nil {
		panic(err)
	}

	medicineRepo, err := database.NewMedicineRepo(db)
	if err != nil {
		panic(err)
	}

	exerciseRepo, err := database.NewExerciseRepo(db)
	if err != nil {
		panic(err)
	}

	apiContext = &api.Context{tremorRepo, medicineRepo, exerciseRepo}

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func serve(request *http.Request) *httptest.ResponseRecorder {
	apiserver := api.NewRouter(apiContext)
	//handler := handlers.LoggingHandler(os.Stdout, apiserver)
	recorder := httptest.NewRecorder()
	//handler.ServeHTTP(recorder, request)
	apiserver.ServeHTTP(recorder, request)
	return recorder
}

func TestPostTremor(t *testing.T) {
	now := time.Now()
	for i := 0; i < 365; i++ {
		resting := 20 + rand.Intn(60)
		postural := 20 + rand.Intn(60)
		date := now.AddDate(0, 0, -1*i)
		tremorJson := fmt.Sprintf(`{"resting": %v, "postural": %v, "date": "%v"}"`,
			resting, postural, date.Format(time.RFC3339))

		request, err := http.NewRequest("POST", "/api/tremors", strings.NewReader(tremorJson))
		if err != nil {
			t.Fatal(err)
		}

		response := serve(request)

		if response.Code != http.StatusOK {
			t.Fatal("Server Error: Returned", response.Code, "instead of", http.StatusOK)
		}
	}
}

func TestGetAllTremors(t *testing.T) {
	request, err := http.NewRequest("GET", "/api/tremors", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := serve(request)

	if response.Code != http.StatusOK {
		t.Fatal("Server Error: Returned", response.Code, "instead of", http.StatusOK)
	}
}

func TestGetTremorsSince(t *testing.T) {
	now := time.Now()

	// test getting tremors for the past week
	then := now.AddDate(0, 0, -6)
	url := "/api/tremors?since=" + then.Format(time.RFC3339)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	response := serve(request)
	if response.Code != http.StatusOK {
		t.Fatal("Server Error: Returned", response.Code, "instead of", http.StatusOK)
	}
	var tremors []api.Tremor
	if err := json.NewDecoder(response.Body).Decode(&tremors); err != nil {
		t.Fatal("decode error:", err)
	}
	for _, tremor := range tremors {
		if tremor.Date.Before(then) {
			t.Fatal("GET request on", url, "returned tremor with timestamp before",
				then.Format(time.RFC3339))
		}
	}
}
