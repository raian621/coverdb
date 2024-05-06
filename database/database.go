package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func CreateDB(filepath, schemaFilepath string) error {
	var err error
	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}

	schemaInitSQL, err := os.ReadFile(schemaFilepath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schemaInitSQL))
	if err != nil {
		return err
	}

	return nil
}
