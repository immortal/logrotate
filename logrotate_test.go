package logrotate

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "TestNew")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	l, err := New(tmpfile.Name(), 0, 0, 1)
	if err != nil {
		t.Error(err)
	}
	if l.Age != 86400 {
		t.Errorf("Expecting age 86400, got: %v", l.Age)
	}
	if l.Num != 7 {
		t.Errorf("Expecting num 7, got: %v", l.Num)
	}
	if l.Size != 1048576 {
		t.Errorf("Expecting size 1048576, got: %v", l.Size)
	}
}
