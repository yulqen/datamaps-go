package reader

import (
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

func TestAlphaStream(t *testing.T) {
	if colstream[26] != "AA" {
		t.Errorf("The test expected AA, got %v.", colstream[26])
	}
	if len(colstream) > maxCols {
		t.Errorf(`Number of columns in alphastream exceeds Excel maximum.
		alphastream contains %d, maxCols is %d`, len(colstream), maxCols)
	}
}

func TestAlphaSingle(t *testing.T) {
	ab := alphabet()
	if ab[0] != "A" {
		t.Errorf("The test expected A, got %v.", ab[0])
	}
	if ab[1] != "B" {
		t.Errorf("The test expected B, got %v.", ab[1])
	}
	if ab[25] != "Z" {
		t.Errorf("The test expected Z, got %v.", ab[25])
	}
}

func TestAlphas(t *testing.T) {
	a := 2 // two alphabets long
	ecs := cols(a)
	cases := []struct {
		col int
		val string
	}{
		{0, "A"},
		{25, "Z"},
		{26, "AA"},
		{52, "BA"},
	}
	for _, c := range cases {
		// we're making sure we can pass that index
		r := 26 * a
		if c.col > r {
			t.Fatalf("Cannot use %d as index to array of %d", c.col, r)
		}
		if got := ecs[c.col]; got != c.val {
			t.Errorf("The test expected ecs[%d] to be %s - got %s.",
				c.col, c.val, ecs[c.col])
		}
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
	_, err := GetTargetFiles(path)
	if err != nil {
		t.Error(err)
	}
}
