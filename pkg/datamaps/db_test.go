package datamaps

import (
	"os"
	"testing"
)

func TestOpenSQLiteFile(t *testing.T) {
	db, err := setupDB("./testdata/test.db")
	defer func() {
		db.Close()
		os.Remove("./testdata/test.db")
	}()

	if err != nil {
		t.Fatalf("%v\ndatamaps-log: Expected to be able to set up the database.", err)
	}

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
	db, err := setupDB("./testdata/test.db")
	defer func() {
		db.Close()
		os.Remove("./testdata/test.db")
	}()

	if err != nil {
		t.Fatal("Expected to be able to set up the database.")
	}

	opts := Options{
		DBPath: "./testdata/test.db",
		DMName: "First Datamap",
		DMPath: "./testdata/datamap.csv",
	}

	if err := DatamapToDB(&opts); err != nil {
		t.Errorf("Unable to write datamap to database file because %v.", err)
	}
}
