package logrotate

import (
	"bytes"
	"os"
	"sync"
	"time"
)

type Logrotate struct {
	sync.Mutex
	Filename string
	Size     int
	Age      int
	Num      int
	file     *os.File
}

func New(logfile string) (*Logrotate, error) {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	return &Logrotate{
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
