package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
	"strconv"
	"time"
)

const (
	medicinesCreate = `create table if not exists medicines(
		mid INTEGER PRIMARY KEY,
		uid INTEGER,
		name TEXT NOT NULL,
		dosage TEXT NOT NULL,
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
	medicineInsert = `insert into medicines(
		uid,
		name,
		dosage,
		mo,	tu, we, th, fr, sa, su,
		reminder,
		startdate,
		enddate)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	medicineSelectBase = "select * from medicines where uid = ?"
	orderByStartDate   = " order by datetime(startdate) desc"
	medicineSelectAll  = medicineSelectBase + orderByStartDate
	medicineSelectMid  = medicineSelectBase + " and mid = ?"
	medicineUpdate     = `update medicines set
		name = ?,
		dosage = ?,
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
		where uid = ? and mid = ?`
	selectForDate = ` and datetime(startdate) < datetime(?2) and
		(enddate is null or datetime(enddate) > datetime(?2))`
	medicineSelectForDate = medicineSelectBase + selectForDate
)

type medicineRepo struct {
	add        *sqlx.Stmt
	getAll     *sqlx.Stmt
	get        *sqlx.Stmt
	getForDate *sqlx.Stmt
	update     *sqlx.Stmt
}

func NewMedicineRepo(db *sqlx.DB) (m *medicineRepo, err error) {
	if _, err = db.Exec(medicinesCreate); err != nil {
		return
	}
	m = new(medicineRepo)
	if m.add, err = db.Preparex(medicineInsert); err != nil {
		return
	}
	if m.getAll, err = db.Preparex(medicineSelectAll); err != nil {
		return
	}
	if m.get, err = db.Preparex(medicineSelectMid); err != nil {
		return
	}
	if m.getForDate, err = db.Preparex(medicineSelectForDate); err != nil {
		return
	}
	if m.update, err = db.Preparex(medicineUpdate); err != nil {
		return
	}
	return
}

func (m *medicineRepo) Add(uid int64, medicine *api.Medicine) (mid int64, err error) {
	if medicine.StartDate == (time.Time{}) {
		medicine.StartDate = time.Now()
	}
	result, err := m.add.Exec(uid,
		medicine.Name,
		medicine.Dosage,
		medicine.Schedule.Mo,
		medicine.Schedule.Tu,
		medicine.Schedule.We,
		medicine.Schedule.Th,
		medicine.Schedule.Fr,
		medicine.Schedule.Sa,
		medicine.Schedule.Su,
		medicine.Reminder,
		medicine.StartDate,
		medicine.EndDate)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (m *medicineRepo) GetAll(uid int64) (medicines []api.Medicine, err error) {
	err = m.getAll.Select(&medicines, uid)
	return
}

func (m *medicineRepo) Get(uid int64, mid int64) (medicine api.Medicine, err error) {
	var medicines []api.Medicine
	if err = m.get.Select(&medicines, uid, mid); err != nil {
		return
	}
	if len(medicines) == 0 {
		err = errors.New("no medicines with MID " + strconv.FormatInt(mid, 10))
		return
	}
	if len(medicines) > 1 {
		err = errors.New("multiple medicines with MID " + strconv.FormatInt(mid, 10))
		return
	}
	medicine = medicines[0]
	return
}

// Returns all medicines scheduled for date (startdate < date < enddate and weekday matches)
func (m *medicineRepo) GetForDate(uid int64, date time.Time) ([]api.Medicine, error) {
	var medicines []api.Medicine
	if err := m.getForDate.Select(&medicines, uid, date); err != nil {
		return nil, err
	}

	weekday := date.Weekday()
	check := func(m api.Medicine) bool {
		switch weekday {
		case time.Monday:
			return m.Schedule.Mo
		case time.Tuesday:
			return m.Schedule.Tu
		case time.Wednesday:
			return m.Schedule.We
		case time.Thursday:
			return m.Schedule.Th
		case time.Friday:
			return m.Schedule.Fr
		case time.Saturday:
			return m.Schedule.Sa
		case time.Sunday:
			return m.Schedule.Su
		}
		return false
	}

	// filter without allocating
	filtered := medicines[:0]
	for _, m := range medicines {
		if check(m) {
			filtered = append(filtered, m)
		}
	}
	return filtered, nil
}

func (m *medicineRepo) Update(uid int64, medicine *api.Medicine) error {
	result, err := m.update.Exec(medicine.Name,
		medicine.Dosage,
		medicine.Schedule.Mo,
		medicine.Schedule.Tu,
		medicine.Schedule.We,
		medicine.Schedule.Th,
		medicine.Schedule.Fr,
		medicine.Schedule.Sa,
		medicine.Schedule.Su,
		medicine.Reminder,
		medicine.StartDate,
		medicine.EndDate,
		uid,
		medicine.MID)
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
