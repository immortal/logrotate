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

// Logrotate struct
type Logrotate struct {
	sync.Mutex
	Age       time.Duration
	Num       int
	Size      int
	file      *os.File
	sTime     time.Time
	size      int64
	timestamp bool
}

// New return instance of Logrotate
// defaults
// age  86400 rotate every 24h0m0s
// num  7     files
// size 0     no limit size
// timestamp false
func New(logfile string, age, num, size int, timestamp bool) (*Logrotate, error) {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	Age := time.Duration(0)
	if age > 0 {
		Age = time.Duration(age) * time.Second
	}
	num--
	if num < 0 {
		num = 7
	}
	Size := 0
	if size > 0 {
		Size = size * 1048576
	}
	lg := &Logrotate{
		Age:       Age,
		Num:       num,
		Size:      Size,
		file:      f,
		sTime:     time.Now(),
		timestamp: timestamp,
	}
	// rotate if needed
	if i, err := lg.file.Stat(); err == nil {
		if lg.Age > 0 && time.Since(i.ModTime()) >= lg.Age {
			if err := lg.rotate(); err != nil {
				return nil, err
			}
		} else if lg.Size > 0 && i.Size() > int64(lg.Size) {
			if err := lg.rotate(); err != nil {
				return nil, err
			}
		}
	}
	return lg, nil
}

// Write implements io.Writer
func (l *Logrotate) Write(p []byte) (n int, err error) {
	l.Lock()
	defer l.Unlock()

	var log []byte

	if l.timestamp {
		t := []byte(time.Now().UTC().Format(time.RFC3339Nano))
		c := [][]byte{t, p}
		log = bytes.Join(c, []byte(" "))
	} else {
		log = p
	}

	writeLen := int64(len(log))

	// rotate based on Age and size
	if l.Age > 0 && time.Since(l.sTime) >= l.Age {
		l.sTime = time.Now()
		if err := l.rotate(); err != nil {
			return 0, err
		}
	} else if l.Size > 0 && l.size+writeLen > int64(l.Size) {
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

// Rotate helper function for rotate
func (l *Logrotate) Rotate() error {
	l.Lock()
	defer l.Unlock()
	return l.rotate()
}

// rotate close existing log file and create a new one
func (l *Logrotate) rotate() error {
	name := l.file.Name()
	l.close()
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
