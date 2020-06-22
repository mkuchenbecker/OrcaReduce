package orcareduce

// Logger is an abstraction over a logging library. Its simple as to remain unopinionated on
// the logging mechanism.
type Logger interface {
	// Errorf logs the format string and args as severity ERROR.
	Errorf(format string, args ...interface{})
	// Error logs an error and message.
	Error(err error, msg string)
	// Infof logs a format string and args as severity INFO.
	Infof(format string, args ...interface{})
}
