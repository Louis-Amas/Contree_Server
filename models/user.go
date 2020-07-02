package models

import (
	"errors"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

// User struct
type User struct {
	ID       int64         `db:"id"`
	Name     string        `json:"name" db:"name"`
	Password string        `json:"-" db:"password"`
	Secret   string        `json:"secret" db:"secret"`
	Socket   socketio.Conn `json:"-"`
}

// InsertUser insert user in db
func (user *User) InsertUser() error {
	res := getDB().MustExec(`INSERT INTO 
	users(name, password, secret) 
	values ($1, $2, $3)
	RETURNING id;

	`, user.Name, user.Password, user.Secret)
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func newUser(name string, password string) *User {
	id := uuid.New().String()
	return &User{
		Name:     name,
		Password: HashString(password),
		Secret:   id,
	}
}

// CheckIfUserExistsAndHasGoodPassword do as he said
func CheckIfUserExistsAndHasGoodPassword(username, password string) (*User, error) {
	var u User
	err := db.Get(&u, "SELECT * FROM users where name=$1 and password=$2", username, HashString(password))
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// VerifyUser return true if user exists and has good secret
func (user *User) VerifyUser() bool {
	var u User
	err := db.Get(&u, "SELECT * FROM users WHERE name=$1", user.Name)
	if err != nil {
		return false
	}
	if user.Secret == u.Secret {
		user.ID = u.ID
		return true
	}
	return false
}

// GetUserFromName return a user from name if user exists otherwise error
func GetUserFromName(name string) (*User, error) {
	var u User
	err := db.Get(&u, "SELECT * FROM users WHERE name=$1", name)
	if err != nil {
		return nil, errors.New("User not found")
	}

	return &u, nil
}
