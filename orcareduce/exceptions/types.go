package exceptions

import (
	"fmt"
)

// PreconditionError is a runtime error that indicates all the preconditions to run a function are not yet satisfied.
// It should generally be used and regarded as a non-serious error that is expected during the normal operation
// of a service or state machine. For example, while its possible to guarentee all preconditions for a function are
// met at the time the function is called it's not possible to guarentee they will still be met at the time of
// preconditions are checked.
type PreconditionError string

// Error implements the error interface function so PreconditionError can be used
// as an error object.
func (err PreconditionError) Error() string {
	return fmt.Sprintf("preconditions not met: %s", string(err))
}

//go:generate mockgen -destination=./mock/mock.go -package=mock github.com/mkuchenbecker/orcareduce/orcareduce/exceptions Logger,Handler
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

type Handler interface {
	// HandleError logs and stats the parameter error and returns that same error to the caller.
	// Supplying a nil error is a no-op. The intent is to be able to have the following:
	//
	// func foo() error {
	//   err := bar()
	//   return handler.HandleError(err)
	// }
	//
	HandleError(error) error

	// HandlePanic captures any panics and populates the parameter err pointer with details about the panic if the
	// error is not already set. HandlePanic should be called via defer.
	//
	// func foo() error {
	//   var err error
	//   done := make(chan bool)
	//   runFunc := func () error {
	//       defer handler.HandlePanic(&err) // called via defer
	//       err := bar() // might panic
	//       done <- true
	//   }
	//   go runFunc()
	//   <- done
	//   return err
	// }
	//
	HandlePanic(err *error)

	// RunAsync safely runs a RunFunc via a goroutine and returns a SyncFunc.
	//
	// runFunc := func () error {
	//     return bar("baz")
	// }
	//
	// syncFunc := handler.RunAsync(runFunc)
	// err := syncFunc()
	//
	//
	RunAsync(f RunFunc) SyncFunc
}

// A SyncFunc is a function with no parameters that returns an error.
type SyncFunc func() error

// A RunFunc is a function with no parameters that returns an error.
type RunFunc func() error
