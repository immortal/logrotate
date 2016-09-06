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

	type Expected struct {
		Age       time.Duration
		Num, Size int
	}

	var testArgs = []struct {
		args     []int
		expected Expected
	}{
		{[]int{0, 0, 0}, Expected{time.Duration(0), 7, 0}},
		{[]int{0, 0, 1}, Expected{time.Duration(0), 7, 1048576}},
		{[]int{0, 1, 0}, Expected{time.Duration(0), 0, 0}},
		{[]int{0, 1, 1}, Expected{time.Duration(0), 0, 1048576}},
		{[]int{1, 0, 0}, Expected{time.Duration(1) * time.Second, 7, 0}},
		{[]int{1, 0, 1}, Expected{time.Duration(1) * time.Second, 7, 1048576}},
		{[]int{1, 1, 1}, Expected{time.Duration(1) * time.Second, 0, 1048576}},
		{[]int{0, 3, 1}, Expected{time.Duration(0), 2, 1048576}},
		{[]int{86400, 0, 1}, Expected{time.Duration(86400) * time.Second, 7, 1048576}},
		{[]int{43200, 0, 1}, Expected{time.Duration(43200) * time.Second, 7, 1048576}},
	}

	for _, a := range testArgs {
		l, err := New(tmpfile.Name(), a.args[0], a.args[1], a.args[2])
		if err != nil {
			t.Error(err)
		}
		if l.Age != a.expected.Age {
			t.Errorf("Expecting age %v, got: %v", a.expected.Age, l.Age)
		}
		if l.Num != a.expected.Num {
			t.Errorf("Expecting num %v, got: %v", a.expected.Num, l.Num)
		}
		if l.Size != a.expected.Size {
			t.Errorf("Expecting size %v, got: %v", a.expected.Size, l.Size)
		}
	}
}

func TestRotate(t *testing.T) {
	var testRotate = []struct {
		args     []int
		expected int
	}{
		{[]int{0, 0, 0}, 1},
		{[]int{0, 0, 1}, 1},
		{[]int{1, 1, 0}, 2},
		{[]int{1, 0, 0}, 4},
		{[]int{1, 3, 0}, 4},
	}

	for _, a := range testRotate {
		dir, err := ioutil.TempDir("", "TestRotate")
		if err != nil {
			t.Error(err)
		}
		tmplog := filepath.Join(dir, "test.log")
		l, err := New(tmplog, a.args[0], a.args[1], a.args[2])
		if err != nil {
			t.Error(err)
		}
		log.SetOutput(l)
		for i := 0; i <= 5; i++ {
			time.Sleep(500 * time.Millisecond)
			log.Println(i)
		}
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != a.expected {
			os.RemoveAll(dir)
			t.Fatalf("Expecting %v got %v", a.expected, len(files))
		}
		os.RemoveAll(dir)
	}
}

func TestRotateIfNotEmpty(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestRotateIfNotEmpty")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)
	tmplog := filepath.Join(dir, "test.log")
	d1 := []byte("not\nempty\n")
	err = ioutil.WriteFile(tmplog, d1, 0644)
	if err != nil {
		t.Error(err)
	}
	l, err := New(tmplog, 0, 0, 0)
	if err != nil {
		t.Error(err)
	}
	log.SetOutput(l)
	for i := 0; i <= 100; i++ {
		log.Println(i)
	}
	fmt.Printf("dir = %+v\n", dir)
	for {

	}

}
