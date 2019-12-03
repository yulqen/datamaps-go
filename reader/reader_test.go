package reader

import (
	"strings"
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

func TestExtract(t *testing.T) {
	d := Extract("testdata/datamap.csv", "testdata/test_template.xlsx")
	cases := []struct {
		sheet, cellref, val string
	}{
		{"Introduction", "C9", "Test Department"},
		{"Introduction", "J9", "Greedy Parrots"},
		{"Introduction", "A1", "10"},
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
	path := "/home/lemon/go/src/github.com/hammerheadlemon/datamaps-go/reader/testdata/"
	files, err := GetTargetFiles(path)
	if err != nil {
		t.Error(err)
	}
	found := false
	for _, f := range files {
		if strings.Contains(f, "test_template.xlsx") {
			found = true
			break
		}
	}
	if found == false {
		t.Errorf("test_template.xlsx not found - and it is there!")
	}
}
