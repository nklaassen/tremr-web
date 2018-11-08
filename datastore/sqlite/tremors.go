package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/datastore"
)

const (
	createTremors = `create table if not exists tremors(
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
	add *sqlx.Stmt
	getAll *sqlx.Stmt
}

func CreateTremorRepo() datastore.TremorRepo {
	db.MustExec(createTremors)
	t := new(tremorRepo)
	var err error
	t.add, err = db.Preparex(tremorInsert)
	if err != nil {
		panic(err)
	}
	t.getAll, err = db.Preparex(tremorSelect)
	if err != nil {
		panic(err)
	}
	return t
}

func (t *tremorRepo) Add(tremor *datastore.Tremor) (err error) {
	_, err = t.add.Exec(tremor.Postural, tremor.Resting, true)
	return
}

func (t *tremorRepo) GetAll() (tremors []datastore.Tremor, err error) {
	err = t.getAll.Select(&tremors)
	return
}
