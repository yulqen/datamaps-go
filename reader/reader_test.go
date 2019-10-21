package reader

import (
	"testing"
)

func TestReadDML(t *testing.T) {
	d, _ := ReadDML("/home/lemon/Documents/datamaps/input/datamap.csv")
	dmlData := *d
	if dmlData[0].Key != "Project/Programme Name" {
		t.Errorf("dmlData[0].key = %s; want Project/Programme Name", dmlData[0].Key)
	}
}
