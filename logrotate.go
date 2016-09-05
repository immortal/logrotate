package logrotate

import (
	"bytes"
	"fmt"
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
	size int64
	file *os.File
}

// New return instance of Logrotate
// defaults
// age  86400 rotate every day
// num  7     files
// size 0     no limit size
func New(logfile string, age, num, size int) (*Logrotate, error) {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
	} else {
		// to test
		size = 1024 * 1024
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
	log := bytes.Join(c, []byte(" "))

	writeLen := int64(len(log))

	fmt.Printf("l.size+writeLen = %+v\n", l.size+writeLen)
	fmt.Printf("l.Size = %+v\n", l.Size)
	if l.size+writeLen > int64(l.Size) {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = l.file.Write(log)
	l.size += int64(n)
	return n, err
}

// Close implements io.Closer, and closes the current logfile
func (l *Logrotate) Close() error {
	l.Lock()
	defer l.Unlock()
	return l.close()
}

// close closes the file if it is open
func (l *Logrotate) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

// Rotate close existing log file and create a new one
func (l *Logrotate) Rotate() error {
	l.Lock()
	defer l.Unlock()
	return l.rotate()
}

func (l *Logrotate) rotate() error {
	if err := l.close(); err != nil {
		return err
	}
	if err := l.openNew(); err != nil {
		return err
	}
	return l.cleanup()
}

func (l *Logrotate) openNew() error {
	name := l.file.Name()
	// rotate logs
	for i := l.Num; i >= 0; i-- {
		logfile := fmt.Sprintf("%s.%d", name, i)
		if _, err := os.Stat(logfile); err == nil {
			// delete old file
			if i == l.Num {
				os.Remove(logfile)
			} else if err := os.Rename(logfile, fmt.Sprintf("%s.%d", name, i+1)); err != nil {
				return err
			}
		}
	}
	// create logfile.log.0
	if err := os.Rename(name, fmt.Sprintf("%s.0", name)); err != nil {
		return err
	}
	// create new log file
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	l.file = f
	l.size = 0
	return nil
}

func (l *Logrotate) cleanup() error {
	return nil
}
