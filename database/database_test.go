package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nklaassen/tremr-web/api"
	"math/rand"
	"os"
	"testing"
	"time"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	if db, err = sqlx.Open("sqlite3", "test_db.sqlite3?_journal_mode=WAL"); err != nil {
		panic(err)
	}

	_, err = db.Exec(`drop table if exists tremors;
		drop table if exists medicines;
		drop table if exists exercises`)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func TestAddTremor(t *testing.T) {
	tremorRepo, err := NewTremorRepo(db)
	if err != nil {
		t.Errorf("Failed to create TremorRepo")
	}

	now := time.Now()
	var tremor api.Tremor

	for i := 0; i < 365; i++ {
		resting := 20 + rand.Intn(60)
		postural := 20 + rand.Intn(60)
		date := now.AddDate(0, 0, -1*i)

		tremor.Resting = &resting
		tremor.Postural = &postural
		tremor.Date = &date

		err = tremorRepo.Add(&tremor)
		if err != nil {
			t.Errorf("Failed to add tremor")
		}
	}
}

func TestGetAllTremors(t *testing.T) {
	tremorRepo, err := NewTremorRepo(db)
	if err != nil {
		t.Errorf("Failed to create TremorRepo")
	}

	_, err = tremorRepo.GetAll()
	if err != nil {
		t.Errorf("Failed to get tremors")
	}
}

func TestGetTremorsSince(t *testing.T) {
	tremorRepo, err := NewTremorRepo(db)
	if err != nil {
		t.Errorf("Failed to create TremorRepo")
	}

	// get values from the past week
	date := time.Now().AddDate(0, 0, -6).Truncate(24 * time.Hour)
	tremors, err := tremorRepo.GetSince(date)
	if err != nil {
		t.Errorf("Failed to get tremors")
	}

	for _, tremor := range tremors {
		if tremor.Date.Before(date) {
			t.Errorf("GetSince returned tremor from before test date")
			t.Errorf("test date: %v, returned date: %v", date, tremor.Date)
		}
	}

	// try getting all tremors since a future date
	date = time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	tremors, err = tremorRepo.GetSince(date)
	if err != nil {
		t.Errorf("Failed to get tremors")
	}

	if len(tremors) != 0 {
		t.Errorf("returned some values from the future")
	}
}

func TestAddMedicine(t *testing.T) {
	medicineRepo, err := NewMedicineRepo(db)
	if err != nil {
		t.Errorf("Failed to create MedicineRepo")
	}

	var medicine api.Medicine
	name := "test med"
	medicine.Name = &name
	dosage := "10 mL"
	medicine.Dosage = &dosage
	medicine.Schedule = &api.Schedule{Mo: true, We: true, Fr: true}
	err = medicineRepo.Add(&medicine)
	if err != nil {
		t.Errorf("Failed to add medicine")
	}
}

func TestGetAllMedicines(t *testing.T) {
	medicineRepo, err := NewMedicineRepo(db)
	if err != nil {
		t.Errorf("Failed to create MedicineRepo")
	}

	_, err = medicineRepo.GetAll()
	if err != nil {
		t.Errorf("Failed to get medicines")
	}
}

func TestAddExercise(t *testing.T) {
	exerciseRepo, err := NewExerciseRepo(db)
	if err != nil {
		t.Errorf("Failed to create ExerciseRepo")
	}

	var exercise api.Exercise
	name := "test exercise"
	exercise.Name = &name
	unit := "10 reps"
	exercise.Unit = &unit
	exercise.Schedule = &api.Schedule{Mo: true, We: true, Fr: true}
	err = exerciseRepo.Add(&exercise)
	if err != nil {
		t.Errorf("Failed to add exercise")
	}
}

func TestGetAllExercises(t *testing.T) {
	exerciseRepo, err := NewExerciseRepo(db)
	if err != nil {
		t.Errorf("Failed to create ExerciseRepo")
	}

	_, err = exerciseRepo.GetAll()
	if err != nil {
		t.Errorf("Failed to get exercises")
	}
}
