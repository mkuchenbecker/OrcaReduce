package orcareduce

// Injector is an object capable of both injecting latency and errors. It's primary use is to
// artificially increase the error rate and latency during testing to help simulate degraded
// performance.
//
//	func (foo *Foo) DoSomethingReliable() error {
//		if err := foo.ChaosInjector.Error(); err != nil {
//			return err
//		}
//		foo.ChaosInjector.Latency() // Blocks the thread for a period of time.
//		... (remainder of function)
//	}
type Injector interface {
	Latency()
	Error() error
}

// ErrorInjector is an object that can inject errors into a process or function. It's primary
// use is to inject errors into otherwise reliable functions for testing purposes.
type ErrorInjector interface {
	Error() error
}

// LatencyInjector is an object that blocks the thread for a period of time. It's primary
// use is to inject latency into otherwise performant functions.
type LatencyInjector interface {
	Latency()
}
