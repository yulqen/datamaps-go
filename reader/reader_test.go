package reader

import (
	"log"
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

func TestReadXLSX(t *testing.T) {
	data := ReadXLSX("testdata/test_template.xlsx")
	if data[0].colL != "A" {
		t.Errorf("Expected data[0].colL to be A, got %v", data[0].colL)
	}
	if data[0].value != "Date:" {
		t.Errorf("Expected data[0].value to be Date:, got %v", data[0].value)
	}
	if data[0].rowLidx != 2 {
		t.Errorf("Expected data[0].rowLidx to be 2, got %v", data[0].rowLidx)
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
	log.Printf("Length of alphastream: %d", len(colstream))
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

func TestReadXLSXToMap(t *testing.T) {
	d := ReadXLSToMap("testdata/datamap.csv", "testdata/test_template.xlsx")
	_ = d
}
