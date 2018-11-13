package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
	"time"
)

const (
	tremorsCreate = `create table if not exists tremors(
		tid INTEGER PRIMARY KEY AUTOINCREMENT,
		postural INTEGER NOT NULL,
		resting INTEGER NOT NULL,
		date DATETIME NOT NULL
	)`
	tremorInsert = "insert into tremors(postural, resting, date) values(?, ?, ?)"
	tremorSelect = "select * from tremors order by datetime(date)"
)

type tremorRepo struct {
	add    *sqlx.Stmt
	getAll *sqlx.Stmt
}

func NewTremorRepo(db *sqlx.DB) (api.TremorRepo, error) {
	_, err := db.Exec(tremorsCreate)
	if err != nil {
		return nil, err
	}
	t := new(tremorRepo)
	t.add, err = db.Preparex(tremorInsert)
	if err != nil {
		return nil, err
	}
	t.getAll, err = db.Preparex(tremorSelect)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *tremorRepo) Add(tremor *api.Tremor) (err error) {
	if tremor.Date == nil {
		now := time.Now()
		tremor.Date = &now
	}
	_, err = t.add.Exec(tremor.Postural, tremor.Resting, tremor.Date)
	return
}

func (t *tremorRepo) GetAll() (tremors []api.Tremor, err error) {
	err = t.getAll.Select(&tremors)
	return
}

func (t *tremorRepo) GetSince(timestamp time.Time) (tremors []api.Tremor, err error) {
	tremors, err = t.GetAll()
	if err != nil {
		return nil, err
	}
	length := len(tremors)
	i := length - 1
	for ; i >= 0; i-- {
		if !timestamp.Before(*tremors[i].Date) {
			break
		}
	}
	return tremors[i+1:], nil
}
