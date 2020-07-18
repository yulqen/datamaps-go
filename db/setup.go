package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDB(path string) (*sql.DB, error) {
	os.Create(path)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return db, errors.New("Cannot open that damn database file")
	}
	stmt := `drop table if exists datamap;
			 create table datamap(id integer no null primary key, name text);
			`
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
	}

	return db, nil
}
