package loggy

import (
	"log"
	"os"
)

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

func New(prefix string) Logger {
	tag := "[" + prefix + "] "
	return &logger{
		tag: tag,
		ll:  log.New(os.Stderr, "", log.LstdFlags),
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
