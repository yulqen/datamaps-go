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
	_, err := ReadDML("/home/bobbins.csv")
	// if we get no error, something has gone wrong
	if err == nil {
		t.Errorf("Should have thrown error %s", err)
	}
}

func TestBadDMLLine(t *testing.T) {
	_, err := ReadDML("/home/lemon/code/python/bcompiler-engine/tests/resources/datamap_empty_cols.csv")
	if err != nil {
		t.Errorf("This will trigger")
	}
}
