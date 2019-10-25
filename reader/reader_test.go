package reader

import (
	"testing"
)

func TestReadDML(t *testing.T) {
	d, _ := ReadDML("/home/lemon/Documents/datamaps/input/datamap.csv")
	dmlData := *d
	// Test Key values
	if dmlData[0].Key != "Project/Programme Name" {
		t.Errorf("dmlData[0].Key = %s; want Project/Programme Name", dmlData[0].Key)
	}
	if dmlData[1].Key != "Department" {
		t.Errorf("dmlData[1].Key = %s; want Department (without a space)", dmlData[1].Key)
	}
	if dmlData[2].Key != "Delivery Body" {
		t.Errorf("dmlData[2].Key = %s; want Delivery Body (without a space)", dmlData[2].Key)
	}
	// Test Sheet values
	if dmlData[0].Sheet != "Introduction" {
		t.Errorf("dmlData[0].Sheet = %s; want Introduction", dmlData[0].Key)
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
