package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	SetLogLevel(INFO)
	log.Println(`init logger with level`, defaultLevel)
}

var defaultLevel uint = 1

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

func SetLogLevel(level uint) {
	if level < 0 || level > 3 {
		return
	}
	defaultLevel = level
}

type LogConf struct {
	Prefix string
	Flags  int
	OutPut io.Writer
}

var leveMap = map[uint]string{
	DEBUG:   `DEBUG`,
	INFO:    `INFO`,
	WARNING: `WARING`,
	ERROR:   `ERROR`,
}

type logger struct {
	prefix   string
	instance *log.Logger
}

func New(lc LogConf) *logger {
	if lc.OutPut == nil {
		lc.OutPut = os.Stderr
	}
	if lc.Flags == 0 {
		lc.Flags = log.Ldate | log.Ltime
	}
	_log := log.New(lc.OutPut, fmt.Sprintf(`%s:%s `, lc.Prefix, leveMap[defaultLevel]), lc.Flags)
	return &logger{lc.Prefix, _log}

}

func (l *logger) Debug(v ...interface{}) {
	output(DEBUG, l, v...)
}

func (l *logger) Info(v ...interface{}) {
	output(INFO, l, v...)
}

func (l *logger) Warning(v ...interface{}) {
	output(WARNING, l, v...)
}

func (l *logger) Error(v ...interface{}) {
	output(ERROR, l, v...)
}

func output(level uint, l *logger, v ...interface{}) {
	if defaultLevel > level {
		return
	}
	l.instance.SetPrefix(fmt.Sprintf(`%s:%s `, l.prefix, leveMap[level]))
	l.instance.Println(v...)
}
