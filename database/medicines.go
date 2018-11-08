package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
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
		enddate DATETIME
	)`
	medicineInsert = `insert into medicines(
		name,
		dosage,
		mo,	tu, we, th, fr, sa, su,
		reminder,
		startdate,
		enddate)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	medicineSelect = "select * from medicines"
)

type medicineRepo struct {
	add    *sqlx.Stmt
	getAll *sqlx.Stmt
}

func NewMedicineRepo(db *sqlx.DB) (api.MedicineRepo, error) {
	_, err := db.Exec(medicinesCreate)
	if err != nil {
		return nil, err
	}
	t := new(medicineRepo)
	t.add, err = db.Preparex(medicineInsert)
	if err != nil {
		return nil, err
	}
	t.getAll, err = db.Preparex(medicineSelect)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *medicineRepo) Add(medicine *api.Medicine) error {
	if medicine.StartDate == nil {
		now := time.Now()
		medicine.StartDate = &now
	}
	_, err := t.add.Exec(medicine.Name,
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

func (t *medicineRepo) GetAll() (medicines []api.Medicine, err error) {
	err = t.getAll.Select(&medicines)
	return
}
