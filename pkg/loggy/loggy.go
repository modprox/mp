package loggy

import (
	"io/ioutil"
	"log"
	"os"
)

// A Logger is used to write log lines to the consul.
type Logger interface {
	Tracef(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}

type logger struct {
	ll  *log.Logger
	tag string
}

// New creates a Logger which prepends each statement with the specified prefix.
func New(prefix string) Logger {
	tag := "[" + prefix + "] "
	return &logger{
		tag: tag,
		ll:  log.New(os.Stderr, "", log.LstdFlags),
	}
}

// Discard creates a Logger which throws away each statement.
func Discard() Logger {
	return &logger{
		ll: log.New(ioutil.Discard, "", 0),
	}
}

func (l logger) Tracef(format string, a ...interface{}) {
	l.printf("TRACE", format, a...)
}

func (l logger) Infof(format string, a ...interface{}) {
	l.printf("INFO ", format, a...)
}

func (l logger) Warnf(format string, a ...interface{}) {
	l.printf("WARN ", format, a...)
}

func (l logger) Errorf(format string, a ...interface{}) {
	l.printf("ERROR", format, a...)
}

func (l logger) printf(level, format string, a ...interface{}) {
	prefixedFmt := level + " " + l.tag + format
	l.ll.Printf(prefixedFmt, a...)
}
