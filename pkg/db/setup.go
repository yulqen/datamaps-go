package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yulqen/datamaps-go/pkg/reader"
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

//DatamapToDB takes a slice of DatamapLine and writes it to a sqlite3 db file.
func DatamapToDB(data []reader.DatamapLine, dm_name string, dm_path string) error {
	log.Printf("Importing Datamap")
	db, err := SetupDB("/home/lemon/.config/datamaps-go/datamaps.db")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	pragma := "PRAGMA foreign_keys = ON;"
	_, err = db.Exec(pragma)
	if err != nil {
		log.Printf("%q: %s\n", err, pragma)
		return err
	}
	stmt_dm, err := tx.Prepare("INSERT INTO datamap (name, date_created) VALUES(?,?)")
	if err != nil {
		return err
	}
	_, err = stmt_dm.Exec(dm_name, time.Now())

	stmt_dml, err := tx.Prepare("INSERT INTO datamap_line (dm_id, key, sheet, cellref) VALUES(?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmt_dm.Close()
	defer stmt_dml.Close()
	for _, dml := range data {
		_, err = stmt_dml.Exec(1, dml.Key, dml.Sheet, dml.Cellref)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
