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
	ecs := cols(2)
	if ecs[0] != "A" {
		t.Errorf("Expected A, got %v", ecs[0])
	}
	if ecs[25] != "Z" {
		t.Errorf("Expected Z, got %v", ecs[25])
	}
	if ecs[26] != "AA" {
		t.Errorf("Expected AA, got %v", ecs[26])
	}
	if ecs[52] != "BA" {
		t.Errorf("Expected BA, got %v", ecs[52])
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
	d := ReadXLSX("testdata/datamap.csv", "testdata/test_template.xlsx")
	if d["Summary"]["A2"].Value != "Date:" {
		t.Errorf("Expected A2 in Summary sheet to be 'Date:' - instead it is %s", d["Summary"]["A2"].Value)
	}
	if d["Another Sheet"]["F5"].Value != "4.2" {
		t.Errorf("Expected F5 in Another Sheet sheet to be 4.2 - instead it is %s", d["Another Sheet"]["F5"].Value)
	}
	if d["Another Sheet"]["J22"].Value != "18" {
		t.Errorf("Expected J22 in Another Sheet sheet to be 18 - instead it is %s", d["Another Sheet"]["J22"].Value)
	}
}

func TestExtract(t *testing.T) {
	d := Extract("testdata/datamap.csv", "testdata/test_template.xlsx")
	if d["Introduction"]["C9"].Value != "Test Department" {
		t.Errorf("Expected C9 in Introduction sheet to be Test Department - instead it is %s", d["Introduction"]["C9"].Value)
	}
	if d["Introduction"]["J9"].Value != "Greedy Parrots" {
		t.Errorf("Expected J9 in Introduction sheet to be Greedy Parrots - instead it is %s", d["Introduction"]["J9"].Value)
	}
	if d["Introduction"]["A1"].Value != "10" {
		t.Errorf("Expected A1 in Introduction sheet to be 10 - instead it is %s", d["Introduction"]["A1"].Value)
	}
}
