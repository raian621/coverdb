package database_test

import (
	"errors"
	"testing"

	"github.com/raian621/coverdb/database"
)

type UserRegisterTest struct {
	username string
	password string
	wantErr  error
}

func TestPostUserSignup(t *testing.T) {
	err := database.CreateDB(DB_PATH, SCHEMA_PATH)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name  string
		users []UserRegisterTest
	}{
		{
			name: "add user to database",
			users: []UserRegisterTest{
				{
					username: "ryan",
					password: "password",
				},
			},
		},
		{
			name: "add duplicate user to database",
			users: []UserRegisterTest{
				{
					username: "ryan",
					password: "password",
				},
				{
					username: "ryan",
					password: "password",
					wantErr:  database.ErrUsernameExists,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			seen := make(map[string]bool, 0)

			for _, user := range tc.users {
				err := database.TryCreateUser(user.username, user.password)
				seen[user.username] = false
				if !errors.Is(err, user.wantErr) {
					t.Errorf("wanted '%v', got '%v' error\n", user.wantErr, err)
				}
			}

			for _, user := range tc.users {
				if seenUser := seen[user.username]; !seenUser {
					seen[user.username] = true
				} else {
					continue
				}

				db := database.GetDB()

				var passhash string
				row := db.QueryRow("SELECT passhash FROM users WHERE username=$1", user.username)
				err := row.Scan(&passhash)
				if err != nil {
					t.Error(err)
				}

				if !database.VerifyHash(user.password, passhash) {
					t.Errorf("password '%s' did not hash to '%s'\n", user.password, passhash)
				}

				_, err = db.Exec("DELETE FROM users WHERE username=$1", user.username)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestSignInUser(t *testing.T) {
	err := database.CreateDB(DB_PATH, SCHEMA_PATH)
	if err != nil {
		t.Fatal(err)
	}
	username := "ryan"
	password := "password"

	if err := database.TryCreateUser(username, password); err != nil {
		t.Error(err)
	}

	db := database.GetDB()
	sessionId, err := database.SignInUser(username, password)
	if err != nil {
		t.Error(err)
	} else {
		_, err = db.Exec("DELETE FROM user_sessions WHERE sessionid=$1", username)
		if err != nil {
			t.Error(err)
		}
	}

	row := db.QueryRow("SELECT username, user_sessions.sessionid FROM user_sessions INNER JOIN users ON user_sessions.user_id=users.id")
	var dbUsername, dbSessionId string
	if err := row.Scan(&dbUsername, &dbSessionId); err != nil {
		t.Error(err)
	}

	if username != dbUsername {
		t.Errorf("wanted '%s', got '%s' username", username, dbUsername)
	}
	if sessionId != dbSessionId {
		t.Errorf("wanted '%s', got '%s' session id", sessionId, dbSessionId)
	}
	_, err = db.Exec("DELETE FROM user_sessions WHERE sessionid=$1", username)
	if err != nil {
		t.Error(err)
	}

	// test incorrect username error
	_, err = database.SignInUser("not the correct user", password)
	if !errors.Is(err, database.ErrInvalidUsername) {
		t.Errorf("wanted '%v', got '%v' error", database.ErrInvalidUsername, err)
	}

	// test incorrect password error
	_, err = database.SignInUser(username, "not the correct password")
	if !errors.Is(err, database.ErrIncorrectPassword) {
		t.Errorf("wanted '%v', got '%v' error", database.ErrIncorrectPassword, err)
	}

	_, err = db.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		t.Error(err)
	}
}
