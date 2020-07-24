package datamaps

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

// mocking funcs in go https://stackoverflow.com/questions/19167970/mock-functions-in-go

func mockConfigDir() (string, error) {
	return "/tmp/CONFIG", nil
}

func TestDBDetect(t *testing.T) {

	if err := os.Mkdir(filepath.Join("/tmp", "CONFIG"), 0700); err != nil {
		t.Fatal("cannot create temporary directory")
	}

	os.Create(filepath.Join("/tmp", "CONFIG", "datamaps.db"))
	defer func() {
		os.RemoveAll(filepath.Join("/tmp", "CONFIG"))
	}()

	dbpc := NewDBPathChecker(mockConfigDir)
	h := dbpc.check()
	log.SetOutput(os.Stderr)
	t.Logf("h is %v\n", h)
	if h != true {
		t.Error("Not there")
	}
}
