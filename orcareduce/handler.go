package orcareduce

// Handler is an object that handles async tasks and the standard processing of errors accross an application
// and safe running of async tasks in the system.
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
// Calling a SyncFunc blocks the thread until the function is called.
type SyncFunc func() error

// A RunFunc is a function with no parameters that returns an error.
type RunFunc func() error
