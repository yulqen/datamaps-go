package reader

import (
	"testing"
)

func TestReadDML(t *testing.T) {
	d, _ := ReadDML("testdata/datamap.csv")
	// Test Key values
	if d[0].Key != "Project/Programme Name" {
		t.Errorf("d[0].Key = %s; want Project/Programme Name", d[0].Key)
	}
	if d[1].Key != "Department" {
		t.Errorf("d[1].Key = %s; want Department (without a space)", d[1].Key)
	}
	if d[2].Key != "Delivery Body" {
		t.Errorf("d[2].Key = %s; want Delivery Body (without a space)", d[2].Key)
	}
	// Test Sheet values
	if d[0].Sheet != "Introduction" {
		t.Errorf("d[0].Sheet = %s; want Introduction", d[0].Key)
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
		t.Errorf("Expected AA, got %v", colstream[26])
	}
	if len(colstream) > maxCols {
		t.Errorf(`Number of columns in alphastream exceeds Excel maximum.
		alphastream contains %d, maxCols is %d`, len(colstream), maxCols)
	}
}

func TestAlphaSingle(t *testing.T) {
	ab := alphabet()
	if ab[0] != "A" {
		t.Errorf("Expected A, got %v", ab[0])
	}
	if ab[1] != "B" {
		t.Errorf("Expected B, got %v", ab[1])
	}
	if ab[25] != "Z" {
		t.Errorf("Expected Z, got %v", ab[25])
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
			t.Errorf("Expected ecs[%d] to be %s - got %s",
				c.col, c.val, ecs[c.col])
		}
	}
}

func TestGetSheetsFromDM(t *testing.T) {
	slice, _ := ReadDML("testdata/datamap.csv")
	sheetNames := getSheetNames(slice)
	if len(sheetNames) != 14 {
		t.Errorf("Expected 14 sheets in slice, got %d",
			len(sheetNames))
	}
}

func TestReadXLSX(t *testing.T) {
	d := ReadXLSX("testdata/test_template.xlsx")
	cases := []struct {
		sheet, cellref, val string
	}{
		{"Summary", "A2", "Date:"},
		{"Another Sheet", "F5", "4.2"},
		{"Another Sheet", "J22", "18"},
	}
	for _, c := range cases {
		got := d[c.sheet][c.cellref].Value
		if got != c.val {
			t.Errorf("Expected %s in %s sheet to be %s "+
				" - instead it is %s", c.cellref, c.sheet, c.val, d[c.sheet][c.cellref].Value)
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
			t.Errorf("Expected %s in %s sheet to be %s "+
				"- instead it is %s", c.sheet, c.cellref, c.val,
				d[c.sheet][c.cellref].Value)
		}
	}
}
