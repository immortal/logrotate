package logrotate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	if l.Age != 86400*time.Second {
		t.Errorf("Expecting age 86400, got: %v", l.Age)
	}
	if l.Num != 7 {
		t.Errorf("Expecting num 7, got: %v", l.Num)
	}
	if l.Size != 1048576 {
		t.Errorf("Expecting size 1048576, got: %v", l.Size)
	}
}

func TestRotate(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestRotate")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir) // clean up
	fmt.Printf("dir = %+v\n", dir)
	tmplog := filepath.Join(dir, "test.log")
	l, err := New(tmplog, 0, 0, 0)
	if err != nil {
		t.Error(err)
	}
	log.SetOutput(l)
	for i := 0; i <= 1000000; i++ {
		time.Sleep(time.Millisecond)
		log.Println(i)
	}
}
