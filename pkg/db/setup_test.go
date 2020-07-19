package db

import (
	"testing"

	"github.com/yulqen/datamaps-go/reader"
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

// TODO:
// THIS IS WHAT I WANT TO DO WITH A TRANNY
// https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
// In the production code of course.

// tx, err := db.Begin()
// if err != nil {
// 	log.Fatal(err)
// }
// stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
// if err != nil {
// 	log.Fatal(err)
// }
// defer stmt.Close()
// for i := 0; i < 100; i++ {
// 	_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
// tx.Commit()

func TestDatamapGoesIntoDB(t *testing.T) {
	db, err := SetupDB("./testdata/test.db") // we set up, then just open it.
	defer db.Close()
	if err != nil {
		t.Fatal("Expected to be able to set up the database.")
	}

	d, _ := reader.ReadDML("testdata/datamap.csv")
	for _, x := range d {
		t.Log(x.Key)
	}

}
