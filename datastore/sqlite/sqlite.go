package sqlite

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
)

var db = sqlx.MustOpen("sqlite3", "db.sqlite3")
