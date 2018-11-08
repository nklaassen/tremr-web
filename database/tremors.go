package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
)

const (
	tremorsCreate = `create table if not exists tremors(
		tid INTEGER PRIMARY KEY AUTOINCREMENT,
		postural INTEGER NOT NULL,
		resting INTEGER NOT NULL,
		completed BOOL NOT NULL,
		date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	tremorInsert = "insert into tremors(postural, resting, completed) values(?, ?, ?)"
	tremorSelect = "select * from tremors"
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
	_, err = t.add.Exec(tremor.Postural, tremor.Resting, true)
	return
}

func (t *tremorRepo) GetAll() (tremors []api.Tremor, err error) {
	err = t.getAll.Select(&tremors)
	return
}
