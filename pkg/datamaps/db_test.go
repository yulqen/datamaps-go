package datamaps

import (
	"database/sql"
	"os"
	"testing"
)

var opts = Options{
	DBPath: "./testdata/test.db",
	DMName: "First Datamap",
	DMPath: "./testdata/short/datamap_matches_test_template.csv",
}

func dbSetup() (*sql.DB, error) {
	db, err := setupDB("./testdata/test.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func dbTeardown(db *sql.DB) {
	db.Close()
	os.Remove("./testdata/test.db")
}

func TestOpenSQLiteFile(t *testing.T) {

	db, err := dbSetup()
	if err != nil {
		t.Fatal(err)
	}
	defer dbTeardown(db)

	stmt := `insert into datamap(id, name) values(1,'cock')`
	_, err = db.Exec(stmt)

	if err != nil {
		t.Errorf("Cannot add record to db")
	}

	rows, err := db.Query("select name from datamap")

	if err != nil {
		t.Errorf("Cannot run select statement")
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			t.Errorf("Cannot scan resulting row")
		}
	}
}

func TestDatamapGoesIntoDB(t *testing.T) {
	db, err := dbSetup()
	if err != nil {
		t.Fatal(err)
	}
	defer dbTeardown(db)

	if err := DatamapToDB(&opts); err != nil {
		t.Errorf("Unable to write datamap to database file because %v.", err)
	}
}

func TestImportSimpleTemplate(t *testing.T) {
	db, err := setupDB("./testdata/test.db")
	defer func() {
		db.Close()
		os.Remove("./testdata/test.db")
	}()

	if err != nil {
		t.Fatal("Expected to be able to set up the database.")
	}

}
