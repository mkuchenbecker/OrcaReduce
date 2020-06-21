package exceptions

import (
	"github.com/mkuchenbecker/orcareduce/orcareduce/exceptions"
	"github.com/pkg/errors"
)

// workflowHandler is a private implementation of the Handler interface.
// It has basic functionality and an internal logger so it can log
// errors it encountners.
type workflowHandler struct {
	logger exceptions.Logger
}

// HandleError logs an error at the appropriate level and returns
// the error as-is. If err is nil, nothing is logged and nil is
// returned.
func (w workflowHandler) HandleError(err error) (out error) {
	out = err
	if err == nil {
		return nil
	}
	cause := errors.Cause(err)
	switch cause.(type) {
	case exceptions.PreconditionError:
		w.logger.Infof(err.Error())
		return
	default:
		w.logger.Error(err, "encountered an error")
		return
	}
}

// HandlePanic recovers from a panic, logs that the panic occured, and populates
// the err paramerter with details about the panic if err is supplied and not
// already populated.
func (w workflowHandler) HandlePanic(err *error) {
	r := recover()
	if r == nil {
		return
	}
	w.logger.Errorf("encountered a panic: %+v", r)
	if err == nil {
		return
	}
	if *err == nil {
		*err = errors.Errorf("encountered a panic: %+v", r)
	}
}

// RunAsync safely runs a RunFunc ansyncronously. It handles any errors that
// are returned by the RunFunc, and recovers from any panics in the RunFunc.
// The return value SyncFunc is a function that blocks until the RunFunc has
// completed execution.
func (w workflowHandler) RunAsync(f exceptions.RunFunc) exceptions.SyncFunc {
	done := make(chan error)
	go func() {
		var err error
		defer func() {
			done <- err
		}()
		defer w.HandlePanic(&err)
		err = f()
		w.HandleError(err)
	}()
	return func() error {
		return <-done
	}
}

// NewHandler makes a Handler using the supplied logger.
func NewHandler(logger exceptions.Logger) exceptions.Handler {
	return workflowHandler{logger: logger}
}
