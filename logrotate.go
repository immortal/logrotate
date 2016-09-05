package logrotate

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

var _ io.WriteCloser = (*Logrotate)(nil)

type Logrotate struct {
	sync.Mutex
	Age  int
	Num  int
	Size int
	file *os.File
}

// New return instance of Logrotate
// defaults
// age  86400 rotate every day
// num  7     files
// size 0     no limit size
func New(logfile string, age, num, size int) (*Logrotate, error) {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
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
func (l *Logrotate) Write(p []byte) (n int, err error) {
	l.Lock()
	defer l.Unlock()

	t := []byte(time.Now().UTC().Format(time.RFC3339Nano))
	c := [][]byte{t, p}
	return l.file.Write(bytes.Join(c, []byte(" ")))
}

// Close implements io.Closer, and closes the current logfile
func (l *Logrotate) Close() error {
	l.Lock()
	defer l.Unlock()
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

// Rotate close existing log file and create a new one
func (l *Logrotate) Rotate() {
}
