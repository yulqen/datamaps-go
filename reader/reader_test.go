package reader

import "testing"

func TestReadDML(t *testing.T) {
	dmlData := *ReadDML("/home/lemon/Documents/datamaps/input/datamap.csv")
	if dmlData[0].Key != "Project/Programme Name" {
		t.Errorf("dmlData[0].key = %s; want Project/Programme Name", dmlData[0].Key)
	}
}
