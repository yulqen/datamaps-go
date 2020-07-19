package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDB(path string) (*sql.DB, error) {
	stmt_base := `DROP TABLE IF EXISTS datamap;
				  CREATE TABLE datamap(id INTEGER PRIMARY KEY, name TEXT, date_created TEXT);
				  DROP TABLE IF EXISTS datamap_line;

				  CREATE TABLE datamap_line(
					id INTEGER PRIMARY KEY,   
					dm_id INTEGER,            
					key TEXT NOT NULL,        
					sheet TEXT NOT NULL,      
					cellref TEXT,             
					FOREIGN KEY (dm_id)       
					REFERENCES datamap(id) 
					ON DELETE CASCADE      
				  );                        
				 `
	os.Create(path)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return db, errors.New("Cannot open that damn database file")
	}

	// We probably don't need pragma here but we have it for later.
	pragma := "PRAGMA foreign_keys = ON;"
	_, err = db.Exec(pragma)
	if err != nil {
		log.Printf("%q: %s\n", err, pragma)
	}

	_, err = db.Exec(stmt_base)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt_base)
	}

	return db, nil
}
