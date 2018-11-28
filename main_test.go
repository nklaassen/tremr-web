package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nklaassen/tremr-web/api"
	"github.com/nklaassen/tremr-web/database"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var router *mux.Router
var token []byte

func TestMain(m *testing.M) {
	// open raw database
	db, err := sqlx.Open("sqlite3", "test_db.sqlite3?_journal_mode=WAL")
	if err != nil {
		panic(err)
	}

	// clear all tables in case the db file already exists
	_, err = db.Exec(`drop table if exists tremors;
		drop table if exists medicines;
		drop table if exists exercises;
		drop table if exists users`)
	if err != nil {
		panic(err)
	}

	// initialize the datastore
	datastore, err := database.GetDataStore(db)
	if err != nil {
		panic(err)
	}

	// set up the api router
	apiEnv := &api.Env{datastore, make(chan struct{})}
	apiRouter := api.NewRouter(apiEnv)

	// setup the global router which strips the /api prefix before sending to the apiRouter
	router = mux.NewRouter()
	router.PathPrefix("/api").Handler(http.StripPrefix("/api", apiRouter))

	// create a user and get an authenticated token to use globally
	user := map[string]interface{}{"email": "test@tremr.com", "password": "hunter2", "name": "tester 1"}
	body, _ := json.Marshal(user)
	_, err = request("POST", "/api/auth/signup", bytes.NewReader(body), http.StatusOK)
	if err != nil {
		panic(err)
	}
	r, err := request("POST", "/api/auth/signin", bytes.NewReader(body), http.StatusOK)
	if err != nil {
		panic(err)
	}
	token = r.Body.Bytes()

	code := m.Run()
	db.Close()
	os.Exit(code)
}

// helper method for performing http requests
func request(method, url string, body io.Reader, expect int) (r *httptest.ResponseRecorder, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	request.Header.Set("Authorization", string(token))
	r = httptest.NewRecorder()
	router.ServeHTTP(r, request)
	if r.Code != expect {
		err = fmt.Errorf("Server Error: Returned %v instead of %v", r.Code, expect)
		return
	}
	return
}

// helper method to generate fractal pseudo-random tremor data
func fractal(a []int) {
	if len(a) <= 2 {
		return
	}
	mid := len(a) / 2
	a[mid] = (a[0] + a[len(a)-1]) / 2
	spread := len(a) / 10
	if spread < 10 {
		spread = 10
	}
	a[mid] += rand.Intn(spread) - spread/2
	if a[mid] > 100 {
		a[mid] = 200 - a[mid]
	}
	if a[mid] < 0 {
		a[mid] = -1 * a[mid]
	}
	fractal(a[:mid+1])
	fractal(a[mid:])
}

func TestPostTremor(t *testing.T) {
	vals := [750]int{}
	vals[0] = 20
	vals[749] = 80
	fractal(vals[:])
	now := time.Now()
	for i := 0; i < 365; i++ {
		resting := vals[i]
		postural := vals[365+i]
		date := now.AddDate(0, 0, -1*i)
		tremorJson := fmt.Sprintf(`{"resting": %v, "postural": %v, "date": "%v"}"`,
			resting, postural, date.Format(time.RFC3339))

		_, err := request(http.MethodPost, "/api/tremors", strings.NewReader(tremorJson), http.StatusOK)
		if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestGetAllTremors(t *testing.T) {
	_, err := request(http.MethodGet, "/api/tremors", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTremorsSince(t *testing.T) {
	now := time.Now()

	// test getting tremors for the past week
	then := now.AddDate(0, 0, -6)
	url := "/api/tremors?since=" + then.Format(time.RFC3339)
	response, err := request(http.MethodGet, url, nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var tremors []api.Tremor
	if err := json.NewDecoder(response.Body).Decode(&tremors); err != nil {
		t.Fatal("decode error:", err)
	}
	// make sure none of the returned tremors are before the timestamp
	for _, tremor := range tremors {
		if tremor.Date.Before(then) {
			t.Error("GET request on", url, "returned tremor with timestamp before",
				then.Format(time.RFC3339))
		}
	}

	// test getting tremors from the future
	then = now.AddDate(0, 0, 1)
	url = "/api/tremors?since=" + then.Format(time.RFC3339)
	response, err = request(http.MethodGet, url, nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(response.Body).Decode(&tremors); err != nil {
		t.Fatal("decode error:", err)
	}
	if len(tremors) > 0 {
		t.Error("Returned tremors from the future!", tremors)
	}
}

func TestPostMedicine(t *testing.T) {
	goodTests := []string{
		`{"name": "test med 1", "dosage": "20 mL", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-11-15T00:00:00Z"}`,
		`{"name": "test med 2", "dosage": "20 mL", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-08-01T00:00:00Z", "enddate": "2018-11-17T00:00:00Z"}`,
		`{"name": "test med 3", "dosage": "20 mL", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-01-01T00:00:00Z", "enddate": "2018-09-01T00:00:00Z"}`,
	}
	badTests := []string{
		`{"dosage": "10 mL", "schedule": {"mo": true, "we": true}}`,
		`{"name": "bad test med 4", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-11-01T00:00:00Z"}`,
		`{"name": "bad test med 5", "dosage": "20 mL", "startdate": "2018-11-01T00:00:00Z"}`,
		`{"name": "test exercise 1", "unit": "10 reps", "schedule": {"mo": true, "we": true}}`}

	for _, test := range goodTests {
		_, err := request(http.MethodPost, "/api/meds", strings.NewReader(test), http.StatusOK)
		if err != nil {
			t.Error(err)
			continue
		}
	}
	for _, test := range badTests {
		_, err := request(http.MethodPost, "/api/meds", strings.NewReader(test), http.StatusBadRequest)
		if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestGetMedicines(t *testing.T) {
	response, err := request(http.MethodGet, "/api/meds", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var medicines []api.Medicine
	if err := json.NewDecoder(response.Body).Decode(&medicines); err != nil {
		t.Fatal("decode error for returned medicines:", err)
		return
	}
}

func TestGetMedicinesForDate(t *testing.T) {
	datestring := "2018-11-27T00:00:00Z" // a tuesday
	date, _ := time.Parse(time.RFC3339, datestring)
	response, err := request(http.MethodGet, "/api/meds?date="+datestring, nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var medicines []api.Medicine
	if err := json.NewDecoder(response.Body).Decode(&medicines); err != nil {
		t.Fatal("decode error for returned medicines:", err)
		return
	}
	for _, med := range medicines {
		if !med.Tu || date.Before(med.StartDate) || (med.EndDate != nil && date.After(*med.EndDate)) {
			t.Error("getForDate returned a medicine not scheduled for that date")
		}
	}
}

func TestGetMedicine(t *testing.T) {
	response, err := request(http.MethodGet, "/api/meds/1", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var medicine api.Medicine
	if err := json.NewDecoder(response.Body).Decode(&medicine); err != nil {
		t.Fatal("decode error for returned medicine:", err)
	}
}

func TestUpdateMedicine(t *testing.T) {

	// get medicine with mid 1
	response, err := request(http.MethodGet, "/api/meds/1", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var medicine api.Medicine
	if err := json.NewDecoder(response.Body).Decode(&medicine); err != nil {
		t.Fatal("decode error for returned medicine:", err)
	}

	// update medicine
	name := "updated medicine"
	medicine.Name = name
	medJson, _ := json.Marshal(medicine)
	response, err = request(http.MethodPut, "/api/meds/1", bytes.NewReader(medJson), http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	// make sure the wrong url gives us a StatusBadRequest
	response, err = request(http.MethodPut, "/api/meds/3", bytes.NewReader(medJson), http.StatusBadRequest)
	if err != nil {
		t.Fatal(err)
	}

	// make sure it actually updated
	response, err = request(http.MethodGet, "/api/meds/1", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(response.Body).Decode(&medicine); err != nil {
		t.Fatal("decode error for returned medicine:", err)
	}
	if medicine.Name != name {
		t.Error("failed to update medicine name")
	}
}

func TestPostExercise(t *testing.T) {
	goodTests := []string{
		`{"name": "test exercise 1", "unit": "20 reps", "schedule": {"mo": true, "we": true, "fr": true},
			"startdate": "2018-11-14T00:00:00Z"}`,
		`{"name": "test exercise 2", "unit": "15 minutes", "schedule": {"tu": true, "th": true, "sa": true},
			"startdate": "2018-03-14T00:00:00Z", "enddate": "2018-11-16T00:00:00Z"}`,
		`{"name": "test exercise 3", "unit": "2 miles", "schedule": {"su": true},
			"startdate": "2017-12-03T00:00:00Z", "enddate": "2018-06-22T00:00:00Z"}`,
	}
	badTests := []string{
		`{"unit": "10 reps", "schedule": {"mo": true, "we": true}}`,
		`{"name": "bad test exercise 4", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-11-01T00:00:00Z"}`,
		`{"name": "bad test exercise 5", "unit": "20 reps", "startdate": "2018-11-01T00:00:00Z"}`,
		`{"name": "test med 1", "dosage": "10 mL", "schedule": {"mo": true, "we": true}}`}

	for _, test := range goodTests {
		_, err := request(http.MethodPost, "/api/exercises", strings.NewReader(test), http.StatusOK)
		if err != nil {
			t.Error(err)
			continue
		}
	}
	for _, test := range badTests {
		_, err := request(http.MethodPost, "/api/exercises", strings.NewReader(test), http.StatusBadRequest)
		if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestGetExercises(t *testing.T) {
	_, err := request(http.MethodGet, "/api/exercises", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExercise(t *testing.T) {
	response, err := request(http.MethodGet, "/api/exercises/1", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var exercise api.Exercise
	if err := json.NewDecoder(response.Body).Decode(&exercise); err != nil {
		t.Fatal("decode error for returned exercise:", err)
	}
}

func TestGetExercisesForDate(t *testing.T) {
	datestring := "2018-11-27T00:00:00Z" // a tuesday
	date, _ := time.Parse(time.RFC3339, datestring)
	response, err := request(http.MethodGet, "/api/exercises?date="+datestring, nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var exercises []api.Exercise
	if err := json.NewDecoder(response.Body).Decode(&exercises); err != nil {
		t.Fatal("decode error for returned exercises:", err)
		return
	}
	for _, exer := range exercises {
		if !exer.Tu || date.Before(exer.StartDate) || (exer.EndDate != nil && date.After(*exer.EndDate)) {
			t.Error("getForDate returned a exercise not scheduled for that date")
		}
	}
}

func TestUpdateExercise(t *testing.T) {
	// get exercise with eid 1
	response, err := request(http.MethodGet, "/api/exercises/1", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var exercise api.Exercise
	if err := json.NewDecoder(response.Body).Decode(&exercise); err != nil {
		t.Fatal("decode error for returned exercise:", err)
	}

	// update exercise
	name := "updated exercise"
	exercise.Name = name
	medJson, _ := json.Marshal(exercise)
	response, err = request(http.MethodPut, "/api/exercises/1", bytes.NewReader(medJson), http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	// make sure the wrong url gives us a StatusBadRequest
	response, err = request(http.MethodPut, "/api/exercises/3", bytes.NewReader(medJson),
		http.StatusBadRequest)
	if err != nil {
		t.Error(err)
	}

	// make sure exercise actually got updated
	response, err = request(http.MethodGet, "/api/exercises/1", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(response.Body).Decode(&exercise); err != nil {
		t.Fatal("decode error for returned exercise:", err)
	}
	if exercise.Name != name {
		t.Error("failed to update exercise name")
	}
}

func TestAuth(t *testing.T) {
	request := func(method string, url string, token []byte, expect int) error {
		request, err := http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}
		request.Header.Set("Authorization", string(token))
		r := httptest.NewRecorder()
		router.ServeHTTP(r, request)
		if r.Code != expect {
			return fmt.Errorf("Server Error: Returned %v instead of %v", r.Code, expect)
		}
		return nil
	}
	// authenticated request
	if err := request(http.MethodGet, "/api/exercises", token, http.StatusOK); err != nil {
		t.Fatal(err)
	}
	// unauthenticated request
	if err := request(http.MethodGet, "/api/exercises", []byte{}, http.StatusUnauthorized); err != nil {
		t.Fatal(err)
	}
	// signup/signin are implicitly tested in TestMain
}

func TestGetUser(t *testing.T) {
	response, err := request(http.MethodGet, "/api/users/1", nil, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var user api.User
	json.NewDecoder(response.Body).Decode(&user)
	if user.Password != "" {
		t.Error("GET /api/users/1 returned the users password!")
	}
	response, err = request(http.MethodGet, "/api/users/2", nil, http.StatusUnauthorized)
	if err != nil {
		t.Fatal(err)
	}
}
