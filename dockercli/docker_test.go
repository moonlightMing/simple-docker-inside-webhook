package dockercli

import "testing"

func TestCleanNoneTagImage(t *testing.T)  {
	if err := CleanNoneTagImage(); err != nil {
		t.Error(err)
	}
}

func TestRemoveImage(t *testing.T) {
	RemoveImage("2b8541a98c0c", false)
}