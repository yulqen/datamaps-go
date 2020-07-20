package db

import (
	"testing"

	"github.com/yulqen/datamaps-go/pkg/reader"
)

func TestOpenSQLiteFile(t *testing.T) {
	db, err := SetupDB("./testdata/test.db")
	defer db.Close()
	if err != nil {
		t.Fatal("Expected to be able to set up the database.")
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
	d, err := reader.ReadDML("./testdata/datamap.csv")
	if err != nil {
		t.Fatal(err)
	}
	err = DatamapToDB(d, "First Datamap", "./testdata/test.db")
	if err != nil {
		t.Errorf("Unable to write datamap to database file because %v.", err)
	}
}
