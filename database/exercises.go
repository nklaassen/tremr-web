package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
	"strconv"
	"time"
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
		startdate DATETIME NOT NULL,
		enddate DATETIME)`
	exerciseInsert = `insert into exercises(
		name,
		unit,
		mo,	tu, we, th, fr, sa, su,
		reminder,
		startdate,
		enddate)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	exerciseSelectBase = "select * from exercises"
	//orderByStartDate   = " order by datetime(startdate)" defined int exercises.go
	exerciseSelectAll = exerciseSelectBase + orderByStartDate
	exerciseSelectEid = exerciseSelectBase + " where eid = ?"
	exerciseUpdate    = `update exercises set
		name = ?,
		unit = ?,
		mo = ?,
		tu = ?,
		we = ?,
		th = ?,
		fr = ?,
		sa = ?,
		su = ?,
		reminder = ?,
		startdate = ?,
		enddate = ?
		where eid = ?`
)

type exerciseRepo struct {
	add    *sqlx.Stmt
	getAll *sqlx.Stmt
	get    *sqlx.Stmt
	update *sqlx.Stmt
}

func NewExerciseRepo(db *sqlx.DB) (apiExerciseRepo api.ExerciseRepo, err error) {
	if _, err = db.Exec(exercisesCreate); err != nil {
		return
	}
	m := new(exerciseRepo)
	if m.add, err = db.Preparex(exerciseInsert); err != nil {
		return
	}
	if m.getAll, err = db.Preparex(exerciseSelectAll); err != nil {
		return
	}
	if m.get, err = db.Preparex(exerciseSelectEid); err != nil {
		return
	}
	if m.update, err = db.Preparex(exerciseUpdate); err != nil {
		return
	}
	apiExerciseRepo = m
	return
}

func (e *exerciseRepo) Add(exercise *api.Exercise) error {
	if exercise.StartDate == nil {
		now := time.Now()
		exercise.StartDate = &now
	}
	if exercise.Reminder == nil {
		f := false
		exercise.Reminder = &f
	}
	_, err := e.add.Exec(exercise.Name,
		exercise.Unit,
		exercise.Schedule.Mo,
		exercise.Schedule.Tu,
		exercise.Schedule.We,
		exercise.Schedule.Th,
		exercise.Schedule.Fr,
		exercise.Schedule.Sa,
		exercise.Schedule.Su,
		exercise.Reminder,
		exercise.StartDate,
		exercise.EndDate)
	return err
}

func (e *exerciseRepo) GetAll() (exercises []api.Exercise, err error) {
	err = e.getAll.Select(&exercises)
	return
}

func (m *exerciseRepo) Get(eid int64) (exercise api.Exercise, err error) {
	var exercises []api.Exercise
	if err = m.get.Select(&exercises, eid); err != nil {
		return
	}
	if len(exercises) == 0 {
		err = errors.New("no exercises with EID " + strconv.FormatInt(eid, 10))
		return
	}
	if len(exercises) > 1 {
		err = errors.New("multiple exercises with EID " + strconv.FormatInt(eid, 10))
		return
	}
	exercise = exercises[0]
	return
}

func (m *exerciseRepo) Update(exercise *api.Exercise) error {
	result, err := m.update.Exec(exercise.Name,
		exercise.Unit,
		exercise.Schedule.Mo,
		exercise.Schedule.Tu,
		exercise.Schedule.We,
		exercise.Schedule.Th,
		exercise.Schedule.Fr,
		exercise.Schedule.Sa,
		exercise.Schedule.Su,
		exercise.Reminder,
		exercise.StartDate,
		exercise.EndDate,
		exercise.EID)
	if err != nil {
		return err
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numRows != 1 {
		return errors.New("Updated " + strconv.FormatInt(numRows, 10) + " rows, expected 1 row")
	}
	return nil
}
