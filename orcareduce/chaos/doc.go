package chaos

/*
chaos is a package to allow the random insertion of errors and latetency. The purpose of injecting latency and errors
is to build a system with the expectation that any function, no matter how reliable, can become transiently unreliable
and should be tested as such.

NewDefault() is a basic meta-injector that will both inject chaos errors and latency. Use of this package would be:

func (foo *Foo) DoSomethingReliable() error {
	if err := foo.ChaosInjector.Error(); err != nil {
		return err
	}
	foo.ChaosInjector.Latency() // Blocks the thread for a period of time.
	... (remainder of function)
}

*/
