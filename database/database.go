package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db             *sql.DB
	filepath       string = "coverdb.db"
	schemaFilepath string = "database/schema.sql"
)

func init() {
	var err error

	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}

	schemaInitSQL, err := os.ReadFile(schemaFilepath)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(schemaInitSQL))
	if err != nil {
		panic(err)
	}
}
