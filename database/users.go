package database

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"
)

var (
	ErrInvalidUsername   error = errors.New("invalid username given")
	ErrIncorrectPassword error = errors.New("invalid password given")
	ErrUsernameExists    error = errors.New("username already exists")
)

func CreateUser(username, password string) error {
	passhash, err := GenerateHash(password)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (username, passhash) VALUES ($1, $2)", username, passhash)

	return err
}

func UsernameTaken(username string) (bool, error) {
	row := db.QueryRow("SELECT COUNT() FROM users WHERE username=$1", username)
	var count int
	if err := row.Scan(&count); err != nil {
		return true, err
	} else if count != 0 {
		return true, nil
	}

	return false, nil
}

func TryCreateUser(username, password string) error {
	if taken, err := UsernameTaken(username); err != nil {
		return err
	} else if taken {
		return ErrUsernameExists
	}

	return CreateUser(username, password)
}

func CreateAdminUser(username, password string) error {
	if taken, err := UsernameTaken(username); err != nil {
		return err
	} else if taken {
		return nil
	}

	if err := CreateUser(username, password); err != nil {
		log.Printf("unexpected error occured while creating admin user: %v", err)
		return err
	}

	// TODO: add permissions and stuff for admin user

	return nil
}

func SignInUser(username, password string) (string, error) {
	// validate credentials
	var userId int
	var passhash string
	row := db.QueryRow(
		"SELECT id, passhash FROM users WHERE username=$1",
		userId,
		passhash,
	)
	if err := row.Scan(&userId, &passhash); err != nil {
		return "", ErrInvalidUsername
	}

	if !VerifyHash(password, passhash) {
		return "", ErrIncorrectPassword
	}

	// generate a random session id
	genSessionId := func() (string, error) {
		randBytes := make([]byte, 30)
		if _, err := rand.Read(randBytes); err != nil {
			return "", err
		}
		return base64.RawStdEncoding.EncodeToString(randBytes), nil
	}

	// try generated sessionIds at most 3 times
	/*
		It's worth noting here that the probability of session ID collisions will
		increase as more users are added to the database, but it won't matter due to
		the size of the session ID.

		The session ID consists of 30 random bytes, which means the total number of
		session IDs, |S|, is 2^240, or around 10^72. To put this into perspective,
		there is estimated to be 10^78 to 10^82 atoms in the known universe.

		The probability of a session ID collision is a function of the number of users
		in the database:

		P(collision) = N / |S|

		and the probability of C collisions is P(collision)^C. Thus, the probability
		of our randomly generated session ID colliding all three 3 times with, let's
		say, 100 billion sessions in the database, is around

		( 10^11 / 10^72 )^3 = ( 10^-51 )^3 = 10^-153

		A session ID with 30 bytes of data is insanely overkill, but I picked the
		session ID length arbitrarily and can't be bothered to fix it in the database.
	*/
	// in retrospect, I should just change the length of the session ID in the DB
	var sessionId string
	uniqueSessionId := false
	for i := 0; i < 3; i++ {
		sessionId, err := genSessionId()
		if err != nil {
			return "", err
		}

		row := db.QueryRow(
			"SELECT COUNT() FROM user_sessions WHERE sessionid=$1",
			sessionId,
		)

		var count int
		if err := row.Scan(&count); err != nil {
			return "", err
		} else if count == 0 {
			uniqueSessionId = true
			break
		}
	}
	// this will probably never happen :D
	if !uniqueSessionId {
		return "", errors.New("this should never happen")
	}

	expires := time.Now().Add(24 * 7 * time.Hour)
	_, err := db.Exec(
		"INSERT INTO user_sessions (sessionid, user_id, expires) VALUES ($1, $2, $3)",
		sessionId,
		userId,
		expires.Unix(),
	)
	if err != nil {
		return "", err
	}

	return "", nil
}
