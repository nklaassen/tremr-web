package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
	"time"
)

const (
	tremorsCreate = `create table if not exists tremors(
		tid INTEGER PRIMARY KEY AUTOINCREMENT,
		uid INTEGER NOT NULL,
		postural INTEGER NOT NULL,
		resting INTEGER NOT NULL,
		date DATETIME NOT NULL
	)`
	tremorInsert      = "insert into tremors(uid, postural, resting, date) values(?, ?, ?, ?)"
	tremorSelectBase  = "select * from tremors where uid = ?"
	orderByDate       = " order by datetime(date)"
	tremorSelectAll   = tremorSelectBase + orderByDate
	tremorSelectSince = tremorSelectBase + " and datetime(date) > datetime(?)" + orderByDate
)

type tremorRepo struct {
	add      *sqlx.Stmt
	getAll   *sqlx.Stmt
	getSince *sqlx.Stmt
}

func NewTremorRepo(db *sqlx.DB) (*tremorRepo, error) {
	_, err := db.Exec(tremorsCreate)
	if err != nil {
		return nil, err
	}
	t := new(tremorRepo)
	t.add, err = db.Preparex(tremorInsert)
	if err != nil {
		return nil, err
	}
	t.getAll, err = db.Preparex(tremorSelectAll)
	if err != nil {
		return nil, err
	}
	t.getSince, err = db.Preparex(tremorSelectSince)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *tremorRepo) Add(uid int64, tremor *api.Tremor) (err error) {
	if tremor.Date == (time.Time{}) {
		tremor.Date = time.Now()
	}
	_, err = t.add.Exec(uid, tremor.Postural, tremor.Resting, tremor.Date)
	return
}

func (t *tremorRepo) GetAll(uid int64) (tremors []api.Tremor, err error) {
	err = t.getAll.Select(&tremors, uid)
	return
}

func (t *tremorRepo) GetSince(uid int64, timestamp time.Time) (tremors []api.Tremor, err error) {
	err = t.getSince.Select(&tremors, uid, timestamp)
	return
}
