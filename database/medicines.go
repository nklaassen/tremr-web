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
		mid INTEGER PRIMARY KEY AUTOINCREMENT,
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
		name,
		dosage,
		mo,	tu, we, th, fr, sa, su,
		reminder,
		startdate,
		enddate)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	medicineSelectBase = "select * from medicines"
	orderByStartDate   = " order by datetime(startdate)"
	medicineSelectAll  = medicineSelectBase + orderByStartDate
	medicineSelectMid  = medicineSelectBase + " where mid = ?"
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
		where mid = ?`
)

type medicineRepo struct {
	add    *sqlx.Stmt
	getAll *sqlx.Stmt
	get    *sqlx.Stmt
	update *sqlx.Stmt
}

func NewMedicineRepo(db *sqlx.DB) (apiMedicineRepo api.MedicineRepo, err error) {
	if _, err = db.Exec(medicinesCreate); err != nil {
		return
	}
	m := new(medicineRepo)
	if m.add, err = db.Preparex(medicineInsert); err != nil {
		return
	}
	if m.getAll, err = db.Preparex(medicineSelectAll); err != nil {
		return
	}
	if m.get, err = db.Preparex(medicineSelectMid); err != nil {
		return
	}
	if m.update, err = db.Preparex(medicineUpdate); err != nil {
		return
	}
	apiMedicineRepo = m
	return
}

func (m *medicineRepo) Add(medicine *api.Medicine) error {
	if medicine.StartDate == nil {
		now := time.Now()
		medicine.StartDate = &now
	}
	if medicine.Reminder == nil {
		f := false
		medicine.Reminder = &f
	}
	_, err := m.add.Exec(medicine.Name,
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
	return err
}

func (m *medicineRepo) GetAll() (medicines []api.Medicine, err error) {
	err = m.getAll.Select(&medicines)
	return
}

func (m *medicineRepo) Get(mid int64) (medicine api.Medicine, err error) {
	var medicines []api.Medicine
	if err = m.get.Select(&medicines, mid); err != nil {
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

func (m *medicineRepo) Update(medicine *api.Medicine) error {
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
