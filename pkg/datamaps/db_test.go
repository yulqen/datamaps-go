package datamaps

import (
	"database/sql"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var singleTarget string = "./testdata/test_template.xlsm"

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
		t.Errorf("unable to write datamap to database file because %v", err)
	}
}

func TestImportSimpleTemplate(t *testing.T) {

	const sql = `
		SELECT return_data.value
		FROM return_data, datamap_line
		WHERE (return_data.filename='test_template.xlsm' AND
			datamap_line.cellref="C9" AND
			return_data.dml_id=datamap_line.id);`

	db, err := dbSetup()
	if err != nil {
		t.Fatal(err)
	}
	defer dbTeardown(db)

	// We need a datamap in there.
	if err := DatamapToDB(&opts); err != nil {
		t.Fatalf("cannot open %s", opts.DMPath)
	}

	if err := importXLSXtoDB(opts.DMName, "TEST RETURN", singleTarget, db); err != nil {
		t.Fatal(err)
	}
	got, err := exec.Command("sqlite3", opts.DBPath, sql).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%q", string(got))
	got_s := string(got)
	got_s = strings.TrimSuffix(got_s, "\n")
	if strings.Compare(got_s, "Test Department") != 0 {
		t.Errorf("we wanted 'Test Department' but got %s", got_s)
	}
}
