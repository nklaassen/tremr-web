package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/nklaassen/tremr-web/api"
)

const (
	usersCreate = `create table if not exists users(
		uid INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT NOT NULL
	)`
	userInsert          = "insert into users(email, password, name) values(?1, ?2, ?3)"
	userSelectFromUid   = "select * from users where uid = ?"
	userSelectFromEmail = "select * from users where email = ?"
)

type userRepo struct {
	add          *sqlx.Stmt
	getFromUid   *sqlx.Stmt
	getFromEmail *sqlx.Stmt
}

func NewUserRepo(db *sqlx.DB) (*userRepo, error) {
	_, err := db.Exec(usersCreate)
	if err != nil {
		return nil, err
	}
	u := new(userRepo)
	u.add, err = db.Preparex(userInsert)
	if err != nil {
		return nil, err
	}
	u.getFromUid, err = db.Preparex(userSelectFromUid)
	if err != nil {
		return nil, err
	}
	u.getFromEmail, err = db.Preparex(userSelectFromEmail)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *userRepo) Add(user *api.User) error {
	_, err := u.add.Exec(user.Email, user.Password, user.Name)
	if err != nil {
		// a user with the same email already exists
		return api.ErrUserExists(err)
	}
	return nil
}

func (u *userRepo) GetFromUid(uid int64) (user api.User, err error) {
	var users []api.User
	err = u.getFromUid.Select(&users, uid)
	if err != nil {
		return
	}
	user = users[0]
	return
}

func (u *userRepo) GetFromEmail(email string) (api.User, error) {
	var users []api.User
	err := u.getFromEmail.Select(&users, email)
	if len(users) == 0 {
		return api.User{}, api.ErrUserDoesNotExist(err)
	}
	if err != nil {
		return api.User{}, err
	}
	return users[0], nil
}
