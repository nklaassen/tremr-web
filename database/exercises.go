package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
)

const (
	exercisesCreate = `create table if not exists exercises(
		eid INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		unit TEXT NOT NULL,
		mo BOOL NOT NULL,
		tu BOOL NOT NULL,
		we BOOL NOT NULL,
		th BOOL NOT NULL,
		fr BOOL NOT NULL,
		sa BOOL NOT NULL,
		su BOOL NOT NULL,
		reminder BOOL NOT NULL,
		startdate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		enddate DATETIME
	)`
	exerciseInsert = `insert into exercises(
		name,
		unit,
		mo,	tu, we, th, fr, sa, su,
		reminder)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	exerciseSelect = "select * from exercises"
)

type exerciseRepo struct {
	add    *sqlx.Stmt
	getAll *sqlx.Stmt
}

func NewExerciseRepo(db *sqlx.DB) (api.ExerciseRepo, error) {
	_, err := db.Exec(exercisesCreate)
	if err != nil {
		return nil, err
	}
	t := new(exerciseRepo)
	t.add, err = db.Preparex(exerciseInsert)
	if err != nil {
		return nil, err
	}
	t.getAll, err = db.Preparex(exerciseSelect)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *exerciseRepo) Add(exercise *api.Exercise) error {
	_, err := t.add.Exec(exercise.Name,
		exercise.Unit,
		exercise.Schedule.Mo,
		exercise.Schedule.Tu,
		exercise.Schedule.We,
		exercise.Schedule.Th,
		exercise.Schedule.Fr,
		exercise.Schedule.Sa,
		exercise.Schedule.Su,
		exercise.Reminder)
	return err
}

func (t *exerciseRepo) GetAll() (exercises []api.Exercise, err error) {
	err = t.getAll.Select(&exercises)
	return
}
