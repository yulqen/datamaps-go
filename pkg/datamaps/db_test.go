package datamaps

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var singleTarget string = "./testdata/test_template.xlsm"

var opts = Options{
	DBPath:   "./testdata/test.db",
	DMName:   "First Datamap",
	DMPath:   "./testdata/datamap_matches_test_template.csv",
	XLSXPath: "./testdata/",
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

	var tests = []struct {
		sheet   string
		cellref string
		value   string
	}{
		{"Introduction", "A1", "10"},
		{"Introduction", "C9", "Test Department"},
		{"Introduction", "C22", "VUNT"},
		{"Introduction", "J9", "Greedy Parrots"},
		{"Summary", "B3", "This is a string"},
		{"Summary", "B4", "2.2"},
		{"Another Sheet", "N34", "23"},
		{"Another Sheet", "DI15", "Rabbit Helga"},
	}

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
		t.Fatalf("Something wrong: %v", err)
	}

	for _, test := range tests {
		sql := fmt.Sprintf(`SELECT return_data.value FROM return_data, datamap_line 
		WHERE 
			(return_data.filename='test_template.xlsm' 
				AND datamap_line.cellref=%q 
				AND datamap_line.sheet=%q
				AND return_data.dml_id=datamap_line.id);`, test.cellref, test.sheet)

		got, err := exec.Command("sqlite3", opts.DBPath, sql).Output()
		if err != nil {
			t.Fatalf("something wrong %v", err)
		}
		got_s := strings.TrimSuffix(string(got), "\n")
		if strings.Compare(got_s, test.value) != 0 {
			t.Errorf("we wanted %s but got %s", test.value, got_s)
		}
	}
}

func TestImportToDB(t *testing.T) {
	var tests = []struct {
		filename string
		sheet    string
		cellref  string
		value    string
	}{
		{"test_template.xlsm", "Introduction", "A1", "10"},
		{"test_template.xlsm", "Introduction", "C9", "Test Department"},
		{"test_template.xlsm", "Introduction", "C22", "VUNT"},
		{"test_template.xlsm", "Introduction", "J9", "Greedy Parrots"},
		{"test_template.xlsm", "Summary", "B3", "This is a string"},
		{"test_template.xlsm", "Summary", "B4", "2.2"},
		{"test_template.xlsm", "Another Sheet", "N34", "23"},
		{"test_template.xlsm", "Another Sheet", "DI15", "Rabbit Helga"},

		{"test_template.xlsx", "Introduction", "A1", "10"},
		{"test_template.xlsx", "Introduction", "C9", "Test Department"},
		{"test_template.xlsx", "Introduction", "C22", "VUNT"},
		{"test_template.xlsx", "Introduction", "J9", "Greedy Parrots"},
		{"test_template.xlsx", "Summary", "B3", "This is a string"},
		{"test_template.xlsx", "Summary", "B4", "2.2"},
		{"test_template.xlsx", "Another Sheet", "N34", "23"},
		{"test_template.xlsx", "Another Sheet", "DI15", "Rabbit Helga"},

		{"test_template2.xlsx", "Introduction", "A1", "10"},
		{"test_template2.xlsx", "Introduction", "C9", "Test Department"},
		{"test_template2.xlsx", "Introduction", "C22", "VUNT"},
		{"test_template2.xlsx", "Introduction", "J9", "Greedy Parrots"},
		{"test_template2.xlsx", "Summary", "B3", "This is a string"},
		{"test_template2.xlsx", "Summary", "B4", "2.2"},
		{"test_template2.xlsx", "Another Sheet", "N34", "23"},
		{"test_template2.xlsx", "Another Sheet", "DI15", "Rabbit Helga"},

		{"test_template3.xlsx", "Introduction", "A1", "10"},
		{"test_template3.xlsx", "Introduction", "C9", "Test Department"},
		{"test_template3.xlsx", "Introduction", "C22", "VUNT"},
		{"test_template3.xlsx", "Introduction", "J9", "Greedy Parrots"},
		{"test_template3.xlsx", "Summary", "B3", "This is a string"},
		{"test_template3.xlsx", "Summary", "B4", "2.2"},
		{"test_template3.xlsx", "Another Sheet", "N34", "23"},
		{"test_template3.xlsx", "Another Sheet", "DI15",
			"Printers run amok in the land of carnivores when bacchus rings 1009.ff    faiioif  !!!]=-=-1290909"},
	}

	db, err := dbSetup()
	if err != nil {
		t.Fatal(err)
	}

	// We need a datamap in there.
	if err := DatamapToDB(&opts); err != nil {
		t.Fatalf("cannot open %s", opts.DMPath)
	}
	defer dbTeardown(db)

	if err := ImportToDB(&opts); err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		sql := fmt.Sprintf(`SELECT return_data.value FROM return_data, datamap_line 
		WHERE 
			(return_data.filename=%q 
				AND datamap_line.cellref=%q 
				AND datamap_line.sheet=%q
				AND return_data.dml_id=datamap_line.id);`, test.filename, test.cellref, test.sheet)

		got, err := exec.Command("sqlite3", opts.DBPath, sql).Output()
		if err != nil {
			t.Fatalf("something wrong %v", err)
		}
		got_s := strings.TrimSuffix(string(got), "\n")
		if strings.Compare(got_s, test.value) != 0 {
			t.Errorf("we wanted value %q in file %s sheet %s %s but got %s",
				test.value, test.filename, test.sheet, test.cellref, got_s)
		}
	}
}

// TODO:

// USING THE INDEX TO tests STRUCT WE COULD DO ALL THESE IN TEST ABOVE

// Returns useful error messages when querying for stuff not in datamap
// func TestImportSimpleQueryValueNotInDatamap(t *testing.T) {
// 	var tests = []struct {
// 		sheet   string
// 		cellref string
// 		value   string
// 	}{
// {"Summary", "B2", "20/10/19"}, // this is not referenced in datamap
// 	}
// }

// TODO:
// When a date is returned from the spreadsheet it is an integer and needs
// to be handled appropriately.
// func TestValuesReturnedAsDates(t *testing.T) {
// }
