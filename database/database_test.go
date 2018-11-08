package database

import (
	"testing"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
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
		t.Errorf("Failed to add tremor")
	}
}
