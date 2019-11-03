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
	_ = ReadXLSX("testdata/test_template.xlsx") // TODO: remove temp blank
}

func TestAlphaStream(t *testing.T) {
	if alphaStream[26] != "AA" {
		t.Errorf("Expected AA, got %v", alphaStream[26])
	}
	if len(alphaStream) > maxCols {
		t.Errorf(`Number of columns in alphastream exceeds Excel maximum.
		alphastream contains %d, maxCols is %d`, len(alphaStream), maxCols)
	}
	log.Printf("Length of alphastream: %d", len(alphaStream))
}

func TestAlphaSingle(t *testing.T) {
	as := alphaSingle()
	if as[0] != "A" {
		t.Errorf("Expected A, got %v", as[0])
	}
	if as[1] != "B" {
		t.Errorf("Expected B, got %v", as[1])
	}
	if as[25] != "Z" {
		t.Errorf("Expected Z, got %v", as[25])
	}
}

func TestAlphas(t *testing.T) {
	as := alphas(2)
	if as[0] != "A" {
		t.Errorf("Expected A, got %v", as[0])
	}
	if as[25] != "Z" {
		t.Errorf("Expected Z, got %v", as[25])
	}
	if as[26] != "AA" {
		t.Errorf("Expected AA, got %v", as[26])
	}
	if as[52] != "BA" {
		t.Errorf("Expected BA, got %v", as[52])
	}

}
