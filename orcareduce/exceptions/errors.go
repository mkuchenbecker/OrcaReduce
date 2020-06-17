package exceptions

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrActorsDepleted      = fmt.Errorf("actors depleted")
	ErrPreconditionsNotMet = fmt.Errorf("preconditions not met")
)

type PreconditionError string

func (err PreconditionError) Error() string {
	return fmt.Sprintf("preconditions not met: %s", err)
}

type Logger interface {
	Error(err error, msg string)
	Info(err error, msg string)
}

type logger struct {
}

func (l logger) Error(err error, msg string) {
	fmt.Printf("encountered an error: %s: %s", msg, err.Error())
}

type Handler interface {
	HandleError(error) error
	HandlePanic(err *error)
}

type workflowHandler struct {
	logger Logger
}

func (w workflowHandler) HandleError(err error) error {
	if err == nil {
		return nil
	}
	cause := errors.Cause(err)
	switch cause.(type) {
	case PreconditionError:
		w.logger.Info(err, "preconditions not met")
		return err
	default:
		w.logger.Info(err, "encountered unknown error")
		return nil
	}
	return nil
}

func (w workflowHandler) HandlePanic(err *error) {
	r := recover()
	if r == nil {
		return
	}
	_ = w.HandleError(r.(error)) // Ignore result.
	if err == nil {
		return
	}
	if *err == nil {
		*err = errors.Errorf("encountered a panic: %s", (*err).Error())
	}
}

type SyncFunc func() error

func (w workflowHandler) RunAsync(f func()) SyncFunc {
	done := make(chan bool)
	go func() {
		defer w.HandlePanic(nil)
		f()
		done <- true
	}()
	return func() error {
		<-done
		return nil
	}
}
