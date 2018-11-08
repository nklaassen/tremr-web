package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nklaassen/tremr-web/api"
	"testing"
)

var db *sqlx.DB

func init() {
	var err error
	if db, err = sqlx.Open("sqlite3", "test_db.sqlite3"); err != nil {
		panic(err)
	}
}

func TestAddTremor(t *testing.T) {
	tremorRepo, err := NewTremorRepo(db)
	if err != nil {
		t.Errorf("Failed to create TremorRepo")
	}

	var tremor api.Tremor
	tremor.Resting = 43
	tremor.Postural = 67
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
	medicine.Name = "testmed"
	medicine.Dosage = "10 mL"
	medicine.Schedule.Mo = true
	medicine.Schedule.We = true
	medicine.Schedule.Fr = true
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
	exercise.Name = "test exercise"
	exercise.Unit = "10 reps"
	exercise.Schedule.Mo = true
	exercise.Schedule.We = true
	exercise.Schedule.Fr = true
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
