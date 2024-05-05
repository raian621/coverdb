package database

import (
	"errors"
)

var ErrUsernameExists error = errors.New("username already exists")

func CreateUser(username, password string) error {
	// check that username isn't taken
	row := db.QueryRow("SELECT COUNT() FROM users WHERE username=$1", username)
	var count int
	row.Scan(&count)

	if count != 0 {
		return ErrUsernameExists
	}

	passhash, err := GenerateHash(password)
	if err != nil {
		return nil
	}
	_, err = db.Exec("INSERT INTO users (username, passhash) VALUES ($1, $2)", username, passhash)

	return err
}
