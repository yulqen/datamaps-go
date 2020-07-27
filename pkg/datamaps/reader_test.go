package datamaps

import (
	"os"
	"reflect"
	"testing"
)

func TestReadDML(t *testing.T) {
	d, _ := ReadDML("testdata/datamap.csv")
	cases := []struct {
		idx int
		val string
	}{
		{0, "Project/Programme Name"},
		{1, "Department"},
		{2, "Delivery Body"},
	}

	for _, c := range cases {
		if got := d[c.idx].Key; got != c.val {
			t.Errorf("The test expected %s, got %s.", c.val, d[c.idx].Key)
		}
	}
}

func TestNoFileReturnsError(t *testing.T) {
	// this file does not exist
	_, err := ReadDML("/home/bobbins.csv")
	// if we get no error, something has gone wrong

	if err == nil {
		t.Errorf("Should have thrown error %s", err)
	}
}

func TestBadDMLLine(t *testing.T) {
	_, err := ReadDML("/home/lemon/code/python/bcompiler-engine/tests/resources/datamap_empty_cols.csv")

	if err == nil {
		t.Errorf("No error so test failed.")
	}
}

func TestGetSheetsFromDM(t *testing.T) {
	slice, _ := ReadDML("testdata/datamap.csv")
	sheetNames := getSheetNames(slice)

	if len(sheetNames) != 15 {
		t.Errorf("The test expected 14 sheets in slice, got %d.",
			len(sheetNames))
	}
}

func TestReadXLSX(t *testing.T) {
	d := ReadXLSX("testdata/test_template.xlsx")
	cases := []struct {
		sheet, cellref, val string
	}{
		{"Summary", "A2", "Date:"},
		{"Summary", "IG10", "botticelli"},
		{"Another Sheet", "F5", "4.2"},
		{"Another Sheet", "J22", "18"},
	}

	for _, c := range cases {
		got := d[c.sheet][c.cellref].Value
		if got != c.val {
			t.Errorf("The test expected %s in %s sheet to be %s "+
				" - instead it is %s.", c.cellref, c.sheet, c.val, d[c.sheet][c.cellref].Value)
		}
	}
}

// func TestExtractWithDBDatamap(t *testing.T) {
// 	// setup - we need the datamap in the test database
// 	db, err := setupDB("./testdata/test.db")
// defer func() {
// 	db.Close()
// 	os.Remove("./testdata/test.db")
// }()

// 	if err != nil {
// 		t.Fatal("Expected to be able to set up the database.")
// 	}

// 	opts := Options{
// 		DBPath: "./testdata/test.db",
// 		DMName: "First Datamap",
// 		DMPath: "./testdata/datamap.csv",
// 	}

// 	if err := DatamapToDB(&opts); err != nil {
// 		t.Errorf("Unable to write datamap to database file because %v.", err)
// 	}

// 	d := extractDBDatamap("First Datamap", "testdata/test_template.xlsx")
// }

func TestDMLSliceFromDatabase(t *testing.T) {
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
		DBPath: "./testdata/test.db",
		DMName: "First Datamap",
		DMPath: "./testdata/datamap.csv",
	}

	if err := DatamapToDB(&opts); err != nil {
		t.Errorf("Unable to write datamap to database file because %v.", err)
	}

	cases := []struct {
		index int
		key   string
	}{
		{0, "Project/Programme Name"},
		{1, "Department"},
		{2, "Delivery Body"},
		{3, "Stoogge value"},
		{4, "DRRDD - IPA ID Number"},
		{5, "Controls Project ID number"},
		{6, "Jokey Entry"},
		{7, "Parrots Name"},
	}

	data, _ := DatamapFromDB("First Datamap", db)

	for _, c := range cases {
		got := data[c.index].Key
		if got != c.key {
			t.Errorf("The test expected %s but got %s\n", c.key, data[c.index].Key)
		}
	}

	if data[0].Key != "Project/Programme Name" {
		t.Errorf("expected to see Project/Programme Name and got %q\n", data[0])
	}
}

func TestExtractUsingDBDM(t *testing.T) {
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
		DBPath: "./testdata/test.db",
		DMName: "First Datamap",
		DMPath: "./testdata/datamap.csv",
	}

	if err := DatamapToDB(&opts); err != nil {
		t.Errorf("Unable to write datamap to database file because %v.", err)
	}

	d, _ := ExtractDBDatamap("First Datamap", "testdata/test_template.xlsx", db)
	cases := []struct {
		sheet, cellref, val string
	}{
		{"Introduction", "C9", "Test Department"},
		{"Introduction", "J9", "Greedy Parrots"},
		{"Introduction", "A1", "10"},
		{"Introduction", "C22", "VUNT"},
	}

	for _, c := range cases {
		got := d[c.sheet][c.cellref].Value
		if got != c.val {
			t.Errorf("The test expected %s in %s sheet to be %s "+
				"- instead it is %s.", c.sheet, c.cellref, c.val,
				d[c.sheet][c.cellref].Value)
		}
	}

	if d["Another Sheet"]["E26"].Value != "Integer:" {
		t.Errorf("Expected E26 in Another Sheet sheet to be Integer: - instead it is %s", d["Another Sheet"]["E26"].Value)
	}
}

func TestExtract(t *testing.T) {
	d := extract("testdata/datamap.csv", "testdata/test_template.xlsx")
	cases := []struct {
		sheet, cellref, val string
	}{
		{"Introduction", "C9", "Test Department"},
		{"Introduction", "J9", "Greedy Parrots"},
		{"Introduction", "A1", "10"},
		{"Introduction", "C22", "VUNT"},
	}

	for _, c := range cases {
		got := d[c.sheet][c.cellref].Value
		if got != c.val {
			t.Errorf("The test expected %s in %s sheet to be %s "+
				"- instead it is %s.", c.sheet, c.cellref, c.val,
				d[c.sheet][c.cellref].Value)
		}
	}

	if d["Another Sheet"]["E26"].Value != "Integer:" {
		t.Errorf("Expected E26 in Another Sheet sheet to be Integer: - instead it is %s", d["Another Sheet"]["E26"].Value)
	}
}

func TestGetTargetFiles(t *testing.T) {
	// This is not a suitable test for parameterisation, but doing it this way anyway.
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"Get files in testdata",
			args{"./testdata/"},
			[]string{
				"testdata/test_template.xlsm",
				"testdata/test_template.xlsx",
				"testdata/test_template2.xlsx",
				"testdata/test_template3.xlsx",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Log("Running the test")
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTargetFiles(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTargetFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTargetFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
