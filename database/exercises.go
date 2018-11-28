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
		eid INTEGER PRIMARY KEY,
		uid INTEGER,
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
		uid,
		name,
		unit,
		mo,	tu, we, th, fr, sa, su,
		reminder,
		startdate,
		enddate)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	exerciseSelectBase = "select * from exercises where uid = ?"
	//orderByStartDate   = " order by datetime(startdate) desc" defined in medicines.go
	exerciseSelectAll = exerciseSelectBase + orderByStartDate
	exerciseSelectEid = exerciseSelectBase + " and eid = ?"
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
		where uid = ? and eid = ?`
	//selectForDate = ` and datetime(startdate) < datetime(?2) and
	//	(enddate is null or datetime(enddate) > datetime(?2))` defined in medicines.go
	exerciseSelectForDate = exerciseSelectBase + selectForDate
)

type exerciseRepo struct {
	add        *sqlx.Stmt
	getAll     *sqlx.Stmt
	get        *sqlx.Stmt
	getForDate *sqlx.Stmt
	update     *sqlx.Stmt
}

func NewExerciseRepo(db *sqlx.DB) (e *exerciseRepo, err error) {
	if _, err = db.Exec(exercisesCreate); err != nil {
		return
	}
	e = new(exerciseRepo)
	if e.add, err = db.Preparex(exerciseInsert); err != nil {
		return
	}
	if e.getAll, err = db.Preparex(exerciseSelectAll); err != nil {
		return
	}
	if e.get, err = db.Preparex(exerciseSelectEid); err != nil {
		return
	}
	if e.getForDate, err = db.Preparex(exerciseSelectForDate); err != nil {
		return
	}
	if e.update, err = db.Preparex(exerciseUpdate); err != nil {
		return
	}
	return
}

func (e *exerciseRepo) Add(uid int64, exercise *api.Exercise) error {
	if exercise.StartDate == (time.Time{}) {
		exercise.StartDate = time.Now()
	}
	_, err := e.add.Exec(uid,
		exercise.Name,
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

func (e *exerciseRepo) GetAll(uid int64) (exercises []api.Exercise, err error) {
	err = e.getAll.Select(&exercises, uid)
	return
}

func (e *exerciseRepo) Get(uid, eid int64) (exercise api.Exercise, err error) {
	var exercises []api.Exercise
	if err = e.get.Select(&exercises, uid, eid); err != nil {
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

// Returns all exercises scheduled for date (startdate < date < enddate and weekday matches)
func (e *exerciseRepo) GetForDate(uid int64, date time.Time) ([]api.Exercise, error) {
	var exercises []api.Exercise
	if err := e.getForDate.Select(&exercises, uid, date); err != nil {
		return nil, err
	}

	weekday := date.Weekday()
	check := func(e api.Exercise) bool {
		switch weekday {
		case time.Monday:
			return e.Schedule.Mo
		case time.Tuesday:
			return e.Schedule.Tu
		case time.Wednesday:
			return e.Schedule.We
		case time.Thursday:
			return e.Schedule.Th
		case time.Friday:
			return e.Schedule.Fr
		case time.Saturday:
			return e.Schedule.Sa
		case time.Sunday:
			return e.Schedule.Su
		}
		return false
	}

	// filter without allocating
	filtered := exercises[:0]
	for _, e := range exercises {
		if check(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

func (e *exerciseRepo) Update(uid int64, exercise *api.Exercise) error {
	result, err := e.update.Exec(exercise.Name,
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
		uid,
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
