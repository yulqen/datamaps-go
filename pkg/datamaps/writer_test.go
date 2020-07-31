package datamaps

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteMaster(t *testing.T) {

	// setup - we need the datamap in the test database
	db, err := setupDB("./testdata/test.db")
	defer func() {
		db.Close()
		os.Remove("./testdata/test.db")
	}()

	if err != nil {
		t.Fatal("Expected to be able to set up the database.")
	}

	opts := Options{
		DBPath:           "./testdata/test.db",
		DMName:           "First Datamap",
		DMPath:           "./testdata/datamap.csv",
		ReturnName:       "Unnamed Return",
		MasterOutPutPath: "./testdata/",
		XLSXPath:         "./testdata/",
	}

	defer func() {
		os.Remove(filepath.Join(opts.MasterOutPutPath, "master.xlsx"))
	}()

	if err := DatamapToDB(&opts); err != nil {
		t.Errorf("Unable to write datamap to database file because %v.", err)
	}

	if err := ImportToDB(&opts); err != nil {
		t.Fatalf("cannot read test XLSX files needed before exporting to master - %v", err)
	}

	if err := ExportMaster(&opts); err != nil {
		t.Error(err)
	}
}
