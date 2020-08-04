package datamaps

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/tealeg/xlsx/v3"
)

var (
	// filesInMaster is a map of each filename from the header row of a master, mapped to its
	// column number
	filesInMaster = make(map[string]int)
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
		DMPath:           "./testdata/datamap_for_master_test.csv",
		ReturnName:       "Unnamed Return",
		MasterOutPutPath: "./testdata/",
		XLSXPath:         "./testdata/",
	}

	// defer func() {
	// 	os.Remove(filepath.Join(opts.MasterOutPutPath, "master.xlsx"))
	// }()

	if err := DatamapToDB(&opts); err != nil {
		t.Fatalf("Unable to write datamap to database file because %v.", err)
	}

	if err := ImportToDB(&opts); err != nil {
		t.Fatalf("cannot read test XLSX files needed before exporting to master - %v", err)
	}

	if err := CreateMaster(&opts); err != nil {
		t.Error(err)
	}

	var tests = []struct {
		key      string
		filename string
		sheet    string
		cellref  string
		value    string
	}{
		{"A Date", "test_template.xlsx", "Summary", "B2", "20/10/19"},
		{"A String", "test_template.xlsx", "Summary", "B3", "This is a string"},
		{"A String2", "test_template.xlsx", "Summary", "C3", "This is a string"},
		{"A String3", "test_template.xlsx", "Summary", "D3", "This is a string"},
		{"A Float", "test_template.xlsx", "Summary", "B4", "2.2"},
		{"An Integer", "test_template.xlsx", "Summary", "B5", "10"},
		{"A Date 1", "test_template.xlsx", "Another Sheet", "B3", "20/10/19"},
		{"A String 1", "test_template.xlsx", "Another Sheet", "B4", "This is a string"},
		{"A Float 1", "test_template.xlsx", "Another Sheet", "B5", "2.2"},
		{"An Integer 1", "test_template.xlsx", "Another Sheet", "B6", "10"},
		{"A Date 2", "test_template.xlsx", "Another Sheet", "D3", "20/10/19"},
		{"A String 2", "test_template.xlsx", "Another Sheet", "D4", "This is a string"},
		{"A Float 3", "test_template.xlsx", "Another Sheet", "D5", "3.2"},
		{"An Integer 3", "test_template.xlsx", "Another Sheet", "D6", "11"},
		{"A Ten Integer", "test_template.xlsx", "Introduction", "A1", "10"},
		{"A Test String", "test_template.xlsx", "Introduction", "C9", "Test Department"},
		{"A Vunt String", "test_template.xlsx", "Introduction", "C22", "VUNT"},
		{"A Parrot String", "test_template.xlsx", "Introduction", "J9", "Greedy Parrots"},
	}

	// Regular testing of import
	// TODO fix date formatting
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
			t.Errorf("when testing the database, for key %s in file %s we expected %s but got %s", test.key,
				test.filename, test.value, got_s)
		}
	}

	// TODO: Testing master

	// The algorthim for this is as follows:
	// - get all the file names in a set from the header row
	// - go through each row and map the key/value coordinates for each file
	// - when going through each test struct, look up the map created in step
	// above to get the value.

	// Open the master and the target sheet
	master, err := xlsx.OpenFile("./testdata/master.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	sheetName := "Master Data"
	sh, ok := master.Sheet[sheetName]
	if !ok {
		t.Errorf("Sheet named %s does not exist", sheetName)
	}
	defer sh.Close()

	err = sh.ForEachRow(rowVisitorTest)

	for _, tt := range tests {
		got, err := masterLookup(sh, tt.key, tt.filename)
		if err != nil {
			t.Fatal("Problem calling masterLookup()")
		}
		if got != tt.value {
			t.Errorf("when testing the master, we expected value %s for key %s in col %s - got %s",
				tt.value, tt.key, tt.filename, got)
		}
	}
}

func masterLookup(sheet *xlsx.Sheet, key string, filename string) (string, error) {
	var out string
	if err := sheet.ForEachRow(func(r *xlsx.Row) error {
		if r.GetCell(0).Value == key {
			out = r.GetCell(filesInMaster[filename]).Value
			return nil
		}
		return nil
	}); err != nil {
		return "", err
	}
	return out, nil
}

func cellVisitorTest(c *xlsx.Cell) error {
	seen := make(map[string]struct{})
	if _, ok := seen[c.Value]; !ok {
		if c.Value != "" {
			x, _ := c.GetCoordinates()
			filesInMaster[c.Value] = x
		}
		seen[c.Value] = struct{}{}
	}

	return nil
}

func rowVisitorTest(r *xlsx.Row) error {
	// TODO here we want to first find the file names from the header row,
	// then test that all key (from col 0) matches the value.

	if r.GetCoordinate() == 0 {
		r.ForEachCell(cellVisitorTest)
		return nil
	}
	return nil
}
