package logger

import (
	"github.com/golang/glog"
	"github.com/mkuchenbecker/orcareduce/orcareduce"
)

// DefaultLogger fetches a simple default logger.
// The logger it returns uses glog as a logger backend.
func DefaultLogger() orcareduce.Logger {
	return logger{
		errorf: glog.Errorf,
		infof:  glog.Infof,
	}
}

type logger struct {
	errorf func(string, ...interface{})
	infof  func(string, ...interface{})
}

// Error logs the error and message as severity Error.
func (l logger) Error(err error, msg string) {
	l.Errorf("%s: %+v", msg, err)
}

// Infof logs the format string and args as severity Info.
func (l logger) Infof(format string, args ...interface{}) {
	l.infof(format, args...)
}

// Errorf logs the format string and args as severity Error
func (l logger) Errorf(format string, args ...interface{}) {
	l.errorf(format, args...)
}
