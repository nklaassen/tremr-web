package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
)

func GetDataStore(db *sqlx.DB) (ds api.DataStore, err error) {
	ds.TremorRepo, err = NewTremorRepo(db)
	if err != nil {
		return
	}
	ds.MedicineRepo, err = NewMedicineRepo(db)
	if err != nil {
		return
	}
	ds.ExerciseRepo, err = NewExerciseRepo(db)
	if err != nil {
		return
	}
	ds.UserRepo, err = NewUserRepo(db)
	if err != nil {
		return
	}
	return
}
