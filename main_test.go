package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	//	"github.com/gorilla/handlers"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nklaassen/tremr-web/api"
	"github.com/nklaassen/tremr-web/database"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var router *mux.Router
var globalAuthTokens []string

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
		drop table if exists users;
		drop table if exists links;`)
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
	handler := http.StripPrefix("/api", apiRouter)
	//handler = handlers.LoggingHandler(os.Stdout, handler)
	router.PathPrefix("/api").Handler(handler)

	// create some users and get authenticated tokens to use globally
	users := []string{
		`{"email": "test1@tremr.com", "password": "hunter1", "name": "tester 1"}`,
		`{"email": "test2@tremr.com", "password": "hunter2", "name": "tester 2"}`,
	}
	for _, user := range users {
		_, err = request("POST", "/api/auth/signup", strings.NewReader(user), "", http.StatusOK)
		if err != nil {
			panic(err)
		}
		r, err := request("POST", "/api/auth/signin", strings.NewReader(user), "", http.StatusOK)
		if err != nil {
			panic(err)
		}
		globalAuthTokens = append(globalAuthTokens, string(r.Body.Bytes()))
	}

	rand.Seed(0xdeadbeef)

	code := m.Run()
	db.Close()
	os.Exit(code)
}

// helper method for performing http requests
func request(method, url string, body io.Reader, token string, expect int) (r *httptest.ResponseRecorder, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	if token != "" {
		request.Header.Set("Authorization", token)
	}
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
	for _, token := range globalAuthTokens {
		vals := [800]int{}
		vals[0] = rand.Intn(10) + 25
		vals[799] = rand.Intn(10) + 75
		fractal(vals[:])
		now := time.Now()
		for i := 0; i < 365; i++ {
			resting := vals[i]
			postural := vals[400+i]
			date := now.AddDate(0, 0, -1*i)
			tremorJson := fmt.Sprintf(`{"resting": %v, "postural": %v, "date": "%v"}"`,
				resting, postural, date.Format(time.RFC3339))

			_, err := request(http.MethodPost, "/api/tremors", strings.NewReader(tremorJson),
				token, http.StatusOK)
			if err != nil {
				t.Error(err)
				continue
			}
		}
	}
}

func TestGetAllTremors(t *testing.T) {
	_, err := request(http.MethodGet, "/api/tremors", nil, globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTremorsSince(t *testing.T) {
	now := time.Now()

	// test getting tremors for the past week
	then := now.AddDate(0, 0, -6)
	url := "/api/tremors?since=" + then.Format(time.RFC3339)
	response, err := request(http.MethodGet, url, nil, globalAuthTokens[0], http.StatusOK)
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
	response, err = request(http.MethodGet, url, nil, globalAuthTokens[0], http.StatusOK)
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
	tests := map[string]int{
		`{"name": "test med 1", "dosage": "20 mL", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-11-15T00:00:00Z"}`: http.StatusOK,
		`{"name": "test med 2", "dosage": "20 mL", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-08-01T00:00:00Z", "enddate": "2018-11-17T00:00:00Z"}`: http.StatusOK,
		`{"name": "test med 3", "dosage": "20 mL", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-01-01T00:00:00Z", "enddate": "2018-09-01T00:00:00Z"}`: http.StatusOK,
		`{"dosage": "10 mL", "schedule": {"mo": true, "we": true}}`: http.StatusBadRequest,
		`{"name": "bad test med 4", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-11-01T00:00:00Z"}`: http.StatusBadRequest,
		`{"name": "bad test med 5", "dosage": "20 mL", "startdate": "2018-11-01T00:00:00Z"}`:   http.StatusBadRequest,
		`{"name": "test exercise 1", "unit": "10 reps", "schedule": {"mo": true, "we": true}}`: http.StatusBadRequest,
	}

	for _, token := range globalAuthTokens {
		for test, expect := range tests {
			_, err := request(http.MethodPost, "/api/meds", strings.NewReader(test), token, expect)
			if err != nil {
				t.Error(err)
				continue
			}
		}
	}
}

func TestGetMedicines(t *testing.T) {
	response, err := request(http.MethodGet, "/api/meds", nil, globalAuthTokens[0], http.StatusOK)
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
	response, err := request(http.MethodGet, "/api/meds?date="+datestring, nil, globalAuthTokens[0], http.StatusOK)
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
	response, err := request(http.MethodGet, "/api/meds/1", nil, globalAuthTokens[0], http.StatusOK)
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
	response, err := request(http.MethodGet, "/api/meds/1", nil, globalAuthTokens[0], http.StatusOK)
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
	response, err = request(http.MethodPut, "/api/meds/1", bytes.NewReader(medJson), globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	// make sure the wrong url gives us a StatusBadRequest
	response, err = request(http.MethodPut, "/api/meds/3", bytes.NewReader(medJson), globalAuthTokens[0], http.StatusBadRequest)
	if err != nil {
		t.Fatal(err)
	}

	// make sure it actually updated
	response, err = request(http.MethodGet, "/api/meds/1", nil, globalAuthTokens[0], http.StatusOK)
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
	tests := map[string]int{
		`{"name": "test exercise 1", "unit": "20 reps", "schedule": {"mo": true, "we": true, "fr": true},
			"startdate": "2018-11-14T00:00:00Z"}`: http.StatusOK,
		`{"name": "test exercise 2", "unit": "15 minutes", "schedule": {"tu": true, "th": true, "sa": true},
			"startdate": "2018-03-14T00:00:00Z", "enddate": "2018-11-16T00:00:00Z"}`: http.StatusOK,
		`{"name": "test exercise 3", "unit": "2 miles", "schedule": {"su": true},
			"startdate": "2017-12-03T00:00:00Z", "enddate": "2018-06-22T00:00:00Z"}`: http.StatusOK,
		`{"unit": "10 reps", "schedule": {"mo": true, "we": true}}`: http.StatusBadRequest,
		`{"name": "bad test exercise 4", "schedule": {"mo": false, "tu": true, "th": true},
			"startdate": "2018-11-01T00:00:00Z"}`: http.StatusBadRequest,
		`{"name": "bad test exercise 5", "unit": "20 reps", "startdate": "2018-11-01T00:00:00Z"}`: http.StatusBadRequest,
		`{"name": "test med 1", "dosage": "10 mL", "schedule": {"mo": true, "we": true}}`:         http.StatusBadRequest,
	}

	for _, token := range globalAuthTokens {
		for test, expect := range tests {
			_, err := request(http.MethodPost, "/api/exercises", strings.NewReader(test), token, expect)
			if err != nil {
				t.Error(err)
				continue
			}
		}
	}
}

func TestGetExercises(t *testing.T) {
	_, err := request(http.MethodGet, "/api/exercises", nil, globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExercise(t *testing.T) {
	response, err := request(http.MethodGet, "/api/exercises/1", nil, globalAuthTokens[0], http.StatusOK)
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
	response, err := request(http.MethodGet, "/api/exercises?date="+datestring, nil, globalAuthTokens[0], http.StatusOK)
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
	response, err := request(http.MethodGet, "/api/exercises/1", nil, globalAuthTokens[0], http.StatusOK)
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
	response, err = request(http.MethodPut, "/api/exercises/1", bytes.NewReader(medJson), globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	// make sure the wrong url gives us a StatusBadRequest
	response, err = request(http.MethodPut, "/api/exercises/3", bytes.NewReader(medJson), globalAuthTokens[0],
		http.StatusBadRequest)
	if err != nil {
		t.Error(err)
	}

	// make sure exercise actually got updated
	response, err = request(http.MethodGet, "/api/exercises/1", nil, globalAuthTokens[0], http.StatusOK)
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
	// basic signup/signin methods tested in testMain

	// test signup with same email
	if _, err := request(http.MethodPost, "/api/auth/signup", strings.NewReader(`{
		"name": "Tester 2",
		"email": "test2@tremr.com",
		"password": "hunter3"
		}`), "", http.StatusConflict); err != nil {
		t.Error(err)
	}

	// test signin with wrong password
	if _, err := request(http.MethodPost, "/api/auth/signin", strings.NewReader(`{
		"email": "test2@tremr.com",
		"password": "hunter3"
		}`), "", http.StatusUnauthorized); err != nil {
		t.Error(err)
	}

	// authenticated request
	if _, err := request(http.MethodGet, "/api/exercises", nil, globalAuthTokens[0], http.StatusOK); err != nil {
		t.Error(err)
	}
	// unauthenticated request
	if _, err := request(http.MethodGet, "/api/exercises", nil, "", http.StatusUnauthorized); err != nil {
		t.Error(err)
	}
	// invalid security token
	if _, err := request(http.MethodGet, "/api/exercises", nil, "0xdeadbeef", http.StatusUnauthorized); err != nil {
		t.Error(err)
	}
}

func TestGetUser(t *testing.T) {
	response, err := request(http.MethodGet, "/api/users/1", nil, globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Error(err)
	}
	var user api.User
	json.NewDecoder(response.Body).Decode(&user)
	if user.Password != "" {
		t.Error("GET /api/users/1 returned the users password!")
	}
	if _, err = request(http.MethodGet, "/api/users/2", nil, globalAuthTokens[0], http.StatusUnauthorized); err != nil {
		t.Error(err)
	}
}

func TestLinks(t *testing.T) {
	// test add link from user 1 to user 2
	_, err := request(http.MethodPost, "/api/users/links/out",
		strings.NewReader(`{"email":"test2@tremr.com"}`), globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Error(err)
	}

	// test add link from user 2 to user 1
	if _, err := request(http.MethodPost, "/api/users/links/out",
		strings.NewReader(`{"email":"test1@tremr.com"}`), globalAuthTokens[1], http.StatusOK); err != nil {
		t.Error(err)
	}

	// test get user 1 incoming links
	response, err := request(http.MethodGet, "/api/users/links/in", nil, globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Error(err)
	}
	var users []api.UserWithoutPassword
	json.NewDecoder(response.Body).Decode(&users)
	if len(users) != 1 || users[0].Email != "test2@tremr.com" {
		t.Error("failed to get incoming links for user 1")
	}

	// test get user 2 outgoing links
	response, err = request(http.MethodGet, "/api/users/links/out", nil, globalAuthTokens[1], http.StatusOK)
	if err != nil {
		t.Error(err)
	}
	json.NewDecoder(response.Body).Decode(&users)
	if len(users) != 1 || users[0].Email != "test1@tremr.com" {
		t.Error("failed to get outgoing links for user 2")
	}

	// test getting user 2 tremors while logged in as user 1
	var tremors1 []api.Tremor
	var tremors2 []api.Tremor
	response, err = request(http.MethodGet, "/api/tremors?uid=2", nil, globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Error(err)
	}
	json.NewDecoder(response.Body).Decode(&tremors1)
	response, err = request(http.MethodGet, "/api/tremors", nil, globalAuthTokens[1], http.StatusOK)
	if err != nil {
		t.Error(err)
	}
	json.NewDecoder(response.Body).Decode(&tremors2)

	if !reflect.DeepEqual(tremors1, tremors2) {
		t.Error("failed to get user 2's tremors while logged in as user 1")
	}

	// test getting user 2 meds while logged in as user 1
	var meds1 []api.Medicine
	var meds2 []api.Medicine
	response, err = request(http.MethodGet, "/api/meds?uid=2", nil, globalAuthTokens[0], http.StatusOK)
	if err != nil {
		t.Error(err)
	}
	json.NewDecoder(response.Body).Decode(&meds1)
	response, err = request(http.MethodGet, "/api/meds", nil, globalAuthTokens[1], http.StatusOK)
	if err != nil {
		t.Error(err)
	}
	json.NewDecoder(response.Body).Decode(&meds2)

	if !reflect.DeepEqual(meds1, meds2) {
		t.Error("failed to get user 2's meds while logged in as user 1")
	}
}
