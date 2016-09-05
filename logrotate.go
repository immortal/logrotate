package logrotate

import (
	"bytes"
	"os"
	"sync"
	"time"
)

type Logrotate struct {
	sync.Mutex
	Age      int
	Filename string
	Num      int
	Size     int
	file     *os.File
}

// New return instance of Logrotate
func New(logfile string, age, num, size int) (*Logrotate, error) {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	// set defaults
	// age  86400 rotate every day
	// num  7 files
	// size 0 no limit size
	if age == 0 {
		age = 86400
	}
	if num == 0 {
		num = 7
	}
	if size > 0 {
		size = size * 1048576
	}
	return &Logrotate{
		Age:  age,
		Num:  num,
		Size: size,
		file: f,
	}, nil
}

// Write implements io.Writer
func (self *Logrotate) Write(p []byte) (n int, err error) {
	self.Lock()
	defer self.Unlock()

	t := []byte(time.Now().UTC().Format(time.RFC3339Nano))
	c := [][]byte{t, p}
	return self.file.Write(bytes.Join(c, []byte(" ")))
}

// Close implements io.Closer, and closes the current logfile
func (self *Logrotate) Close() error {
	self.Lock()
	defer self.Unlock()
	if self.file == nil {
		return nil
	}
	err := self.file.Close()
	self.file = nil
	return err
}
