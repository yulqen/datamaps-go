package datamaps

import (
	"os"
	"path/filepath"
	"testing"
)

// mocking funcs in go https://stackoverflow.com/questions/19167970/mock-functions-in-go

func mockConfigDir() (string, error) {
	return "/tmp/CONFIG/datamaps/", nil
}

func TestDBDetect(t *testing.T) {

	cPath := filepath.Join("/tmp", "CONFIG", "datamaps")
	t.Logf("%s\n", cPath)
	if err := os.MkdirAll(cPath, 0700); err != nil {
		t.Fatalf("cannot create temporary directory - %v", err)
	}

	os.Create(filepath.Join("/tmp", "CONFIG", "datamaps", "datamaps.db"))
	defer func() {
		os.RemoveAll(filepath.Join("/tmp", "CONFIG"))
	}()

	dbpc := NewDBPathChecker(mockConfigDir)
	h := dbpc.Check()
	if !h {
		t.Error("the db file should be found but isn't")
	}
}
