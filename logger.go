package golog

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	Ldate = 1 << iota
	Ltime
)

type Logger struct {
	mu       sync.Mutex
	options  *Options
	size     int64
	index    int64
	file     *os.File
	fileTime time.Time
	flags    int
}

type Options struct {
	LogDir     string `yaml:"log_dir"`
	MaxSize    int64  `yaml:"max_size"`
	Daily      bool   `yaml:"daily"`
	LogName    string `yaml:"log_name"`
	DateFormat string `yaml:"date_format"`
}

func NewLogger(options *Options, flags int) *Logger {
	logger := &Logger{
		size:    0,
		index:   0,
		options: options,
	}
	logger.flags = flags
	logger.rotate()
	return logger
}

func isCurrentDay(date time.Time) bool {
	y, m, d := date.Date()
	y1, m1, d1 := time.Now().Date()
	return y == y1 && m == m1 && d == d1
}

func (l *Logger) logName() string {
	dateFormat := "01-02-2006"
	if l.options.DateFormat != "" {
		dateFormat = l.options.DateFormat
	}

	if l.index == 0 {
		if l.options.Daily {
			return fmt.Sprintf("%s_%s.log", l.options.LogName, time.Now().Format(dateFormat))
		} else {
			return fmt.Sprintf("%s.log", l.options.LogName)
		}
	} else {
		if l.options.Daily {
			return fmt.Sprintf("%s_%s_%d.log", l.options.LogName, time.Now().Format(dateFormat), l.index)
		} else {
			return fmt.Sprintf("%s_%d.log", l.options.LogName, l.index)
		}
	}
}

func (l *Logger) checkRotate() {
	if !isCurrentDay(l.fileTime) || (l.options.MaxSize > 0 && l.size >= l.options.MaxSize) {
		l.index++
		if l.file != nil {
			l.file.Close()
		}
		l.rotate()
	}
}

func (l *Logger) rotate() {
	os.MkdirAll(l.options.LogDir, os.ModePerm)
	var err error
	if l.file != nil {
		l.file.Close()
	}

	l.file, err = os.OpenFile(l.logName(), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	stat, err := l.file.Stat()
	if err != nil {
		l.file.Close()
		return
	}

	l.size = stat.Size()
	l.fileTime = stat.ModTime()

	l.checkRotate()
}

func (l *Logger) prefix(tag ...string) string {
	prefix := ""

	if l.flags&(Ldate|Ltime) != 0 {
		if l.flags&Ldate != 0 {
			prefix = time.Now().Format("01-02-2006")
		}
		if l.flags&Ltime != 0 {
			prefix += " " + time.Now().Format("15:04:05")
		}
	}

	if len(tag) > 0 {
		return fmt.Sprintf("[%s] %s - ", strings.Trim(prefix, ""), tag[0])
	}

	return fmt.Sprintf("[%s] - ", strings.Trim(prefix, ""))
}

func (l *Logger) write(data []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.checkRotate()

	n, _ := l.file.Write(data)
	l.size += int64(n)
}

func (l *Logger) writeString(data string, tag ...string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.checkRotate()

	n, _ := l.file.WriteString(l.prefix(tag...) + data + "\n")
	l.size += int64(n)
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

func (l *Logger) Println(v string, tag ...string) {
	l.writeString(v, tag...)
}

func (l *Logger) Info(v string) {
	l.writeString(v, "Info")
}

func (l *Logger) Error(v string) {
	l.writeString(v, "Error")
}

func (l *Logger) Warning(v string) {
	l.writeString(v, "Warning")
}
