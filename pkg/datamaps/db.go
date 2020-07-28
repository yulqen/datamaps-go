package datamaps

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	// Needed for the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// setupDB creates the intitial database
func setupDB(path string) (*sql.DB, error) {
	stmtBase := `DROP TABLE IF EXISTS datamap;
				 DROP TABLE IF EXISTS datamap_line;
				 DROP TABLE IF EXISTS return;
				 DROP TABLE IF EXISTS return_data;

				  CREATE TABLE datamap(
					  id INTEGER PRIMARY KEY,
					  name TEXT,
					  date_created TEXT);

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

				 CREATE TABLE return(
					 id INTEGER PRIMARY KEY,
					 name TEXT,
					 date_created TEXT
					);

				 CREATE TABLE return_data(
					 id INTEGER PRIMARY KEY,
					 dml_id INTEGER,
					 ret_id INTEGER,
					 value TEXT,
					 FOREIGN KEY (dml_id)
					 REFERENCES datamap_line(id) 
					 ON DELETE CASCADE
					 FOREIGN KEY (ret_id)
					 REFERENCES return(id) 
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
		// log.Printf("%q: %s\n", err, pragma)
		return nil, err
	}

	_, err = db.Exec(stmtBase)
	if err != nil {
		// log.Printf("%q: %s\n", err, stmt_base)
		return nil, err
	}

	return db, nil
}

// Import a directory of xlsx files to the database, using the datamap
// to filter the data.
func ImportToDB(opts *Options) error {
	fmt.Printf("Import files in %s\n\tas return named %s\n\tusing datamap named %s\n", opts.XLSXPath, opts.ReturnName, opts.DMName)

	target, err := getTargetFiles(opts.XLSXPath)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", opts.DBPath)
	if err != nil {
		return err
	}

	for _, vv := range target {
		// TODO: Do the work!

		if err := importXLSXtoDB(opts.DMName, opts.ReturnName, vv, db); err != nil {
			return err
		}
	}
	return nil
}

// DatamapToDB takes a slice of datamapLine and writes it to a sqlite3 db file.
func DatamapToDB(opts *Options) error {
	fmt.Printf("Importing datamap file %s and naming it %s.\n", opts.DMPath, opts.DMName)

	data, err := ReadDML(opts.DMPath)
	if err != nil {
		log.Fatal(err)
	}

	d, err := sql.Open("sqlite3", opts.DBPath)
	if err != nil {
		return errors.New("Cannot open that damn database file")
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	pragma := "PRAGMA foreign_keys = ON;"
	_, err = d.Exec(pragma)
	if err != nil {
		log.Printf("%q: %s\n", err, pragma)
		return err
	}

	stmtDm, err := tx.Prepare("INSERT INTO datamap (name, date_created) VALUES(?,?)")
	if err != nil {
		return err
	}

	res, err := stmtDm.Exec(opts.DMName, time.Now())
	if err != nil {
		return err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	stmtDml, err := tx.Prepare("INSERT INTO datamap_line (dm_id, key, sheet, cellref) VALUES(?,?,?,?);")
	if err != nil {
		return err
	}

	defer stmtDm.Close()
	defer stmtDml.Close()

	for _, dml := range data {
		_, err = stmtDml.Exec(lastId, dml.Key, dml.Sheet, dml.Cellref)
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

func importXLSXtoDB(dm_name string, return_name string, file string, db *sql.DB) error {
	// d, err := ExtractDBDatamap(dm_name, file, db)
	d, err := ExtractDBDatamap(dm_name, file, db)
	if err != nil {
		return err
	}
	fmt.Printf("Extracting from %s\n", file)
	// fmt.Printf("Data is: %#v\n", d["Introduction"]["C17"].Value)

	stmtReturn, err := db.Prepare("insert into return(name, date_created) values(?,?)")
	if err != nil {
		log.Fatal(err)
	}

	defer stmtReturn.Close()

	res, err := stmtReturn.Exec(return_name, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	retId, err := res.LastInsertId()
	fmt.Println(retId)
	if err != nil {
		log.Fatal(err)
	}

	// We're going to need a transaction for the big stuff
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for sheetName, sheetData := range d {

		for cellRef, cellData := range sheetData {
			// fmt.Printf("Getting %s from sheet %s\n", cellRef, sheetName)

			dmlQuery, err := db.Prepare("select id from datamap_line where (sheet=? and cellref=?)")
			if err != nil {
				log.Fatal(err)
			}
			defer dmlQuery.Close()
			dmlIdRow := dmlQuery.QueryRow(sheetName, cellRef)
			fmt.Println(dmlIdRow)

			var dmlId *int

			if err := dmlIdRow.Scan(&dmlId); err != nil {
				fmt.Errorf("cannot find a datamap_line row for %s and %s: %s\n", sheetName, cellRef, err)
			}

			insertStmt, err := db.Prepare("insert into return_data (dml_id, value) values(?,?)")
			if err != nil {
				log.Fatal(err)
			}
			defer insertStmt.Close()

			res, err = insertStmt.Exec(dmlId, cellData.Value)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
