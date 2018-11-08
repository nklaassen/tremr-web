package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nklaassen/tremr-web/api"
	"testing"
	"math/rand"
	"time"
)

var db *sqlx.DB

func init() {
	var err error
	if db, err = sqlx.Open("sqlite3", "test_db.sqlite3"); err != nil {
		panic(err)
	}
	_, err = db.Exec("drop table tremors; drop table medicines; drop table exercises")
	if err != nil {
		panic(err)
	}
}

func TestAddTremor(t *testing.T) {
	tremorRepo, err := NewTremorRepo(db)
	if err != nil {
		t.Errorf("Failed to create TremorRepo")
	}

	var tremor api.Tremor
	resting := 20 + rand.Intn(60)
	postural := 20 + rand.Intn(60)
	tremor.Resting = &resting
	tremor.Postural = &postural
	now := time.Now()
	tremor.Date = &now

	err = tremorRepo.Add(&tremor)
	if err != nil {
		t.Errorf("Failed to add tremor")
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
