package datamaps

import "testing"

func Test_getUserConfigDir(t *testing.T) {
	dir, _ := getUserConfigDir()
	if dir != "/home/lemon/.config/datamaps-go" {
		t.Errorf("Did not find the correct directory - found %s instead\n", dir)
	}
}
