package datamaps

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
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
					 filename TEXT,
					 value TEXT,
					 FOREIGN KEY (dml_id)
					 REFERENCES datamap_line(id) 
					 ON DELETE CASCADE
					 FOREIGN KEY (ret_id)
					 REFERENCES return(id) 
					 ON DELETE CASCADE
				 );
				 `
	if _, err := os.Create(path); err != nil {
		return nil, err
	}
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

// ImportToDB imports a directory of xlsx files to the database, using the datamap
// to filter the data.
func ImportToDB(opts *Options) error {
	log.Printf("Importing files in %s as return named %s using datamap named %s.", opts.XLSXPath, opts.ReturnName, opts.DMName)

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
	log.Printf("Importing datamap file %s and naming it %s.\n", opts.DMPath, opts.DMName)

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

	lastID, err := res.LastInsertId()
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
		_, err = stmtDml.Exec(lastID, dml.Key, dml.Sheet, dml.Cellref)
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

func importXLSXtoDB(dmName string, returnName string, file string, db *sql.DB) error {
	// d, err := ExtractDBDatamap(dmName, file, db)
	_, filename := path.Split(file)
	d, err := ExtractDBDatamap(dmName, file, db)
	if err != nil {
		return err
	}
	log.Printf("Extracting from %s.\n", file)

	// If there is already a return with a matching name, use that.
	rtnQuery, err := db.Prepare("select id from return where (return.name=?)")
	if err != nil {
		log.Fatal(err)
	}
	defer rtnQuery.Close()

	var retID int64
	err = rtnQuery.QueryRow(returnName).Scan(&retID)
	if err != nil {
		log.Printf("Couldn't find an existing return named '%s' so let's create it.\n", returnName)
	}

	if retID == 0 {
		stmtReturn, err := db.Prepare("insert into return(name, date_created) values(?,?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmtReturn.Close()

		res, err := stmtReturn.Exec(returnName, time.Now())
		if err != nil {
			err := fmt.Errorf("%v\nCannot create %s", err, returnName)
			log.Println(err.Error())
			os.Exit(1)
		}

		retID, err = res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
	}

	// We're going to need a transaction for the big stuff
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for sheetName, sheetData := range d {

		for cellRef, cellData := range sheetData {

			dmlQuery, err := db.Prepare("select id from datamap_line where (sheet=? and cellref=?)")
			if err != nil {
				log.Fatal(err)
			}
			defer dmlQuery.Close()
			dmlIDRow := dmlQuery.QueryRow(sheetName, cellRef)

			var dmlID *int

			if err := dmlIDRow.Scan(&dmlID); err != nil {
				err := fmt.Errorf("cannot find a datamap_line row for %s and %s: %s", sheetName, cellRef, err)
				log.Println(err.Error())
			}

			insertStmt, err := db.Prepare("insert into return_data (dml_id, ret_id, filename, value) values(?,?,?,?)")
			if err != nil {
				log.Fatal(err)
			}
			defer insertStmt.Close()

			_, err = insertStmt.Exec(dmlID, retID, filename, cellData.Value)
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
