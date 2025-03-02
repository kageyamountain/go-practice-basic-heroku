package models

import (
	"log"
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	Todos     []Todo
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

func (u *User) CreateUser() (err error) {
	cmd := `INSERT INTO users (uuid, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = Db.Exec(cmd, createUUID(), u.Name, u.Email, Encrypt(u.Password), time.Now())

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetUser(id int) (user User, err error) {
	user = User{}
	cmd := `SELECT id, uuid, name, email, password, created_at FROM users WHERE id = $1`
	err = Db.QueryRow(cmd, id).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

func (u *User) UpdateUser() (err error) {
	cmd := `
		UPDATE users
		SET
			name = $1,
			email = $2
		WHERE
		    id = $3`

	_, err = Db.Exec(cmd, u.Name, u.Email, u.Id)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func (u *User) DeleteUser() (err error) {
	cmd := `DELETE FROM users WHERE id = $1`

	_, err = Db.Exec(cmd, u.Id)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetUserByEmail(email string) (user User, err error) {
	user = User{}
	cmd := `SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1`
	err = Db.QueryRow(cmd, email).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

func (u *User) CreateSession() (session Session, err error) {
	session = Session{}
	cmd1 := `INSERT INTO sessions (uuid, email, user_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err = Db.Exec(cmd1, createUUID(), u.Email, u.Id, time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	cmd2 := `SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1 AND email = $2`
	err = Db.QueryRow(cmd2, u.Id, u.Email).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)

	return session, err
}

func (session *Session) CheckSession() (valid bool, err error) {
	cmd := `SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1`
	err = Db.QueryRow(cmd, session.Uuid).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return valid, err
	}

	if session.Id != 0 {
		valid = true
	}
	return valid, err
}

func (session *Session) DeleteSessionByUuid() (err error) {
	cmd := `DELETE FROM sessions WHERE uuid = $1`
	_, err = Db.Exec(cmd, session.Uuid)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (session *Session) GetUserBySession() (user User, err error) {
	user = User{}
	cmd := `SELECT id, uuid, name, email, created_at FROM users WHERE id = $1`
	err = Db.QueryRow(cmd, session.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return user, err
}
