package orcareduce

type Injector interface {
	Latency()
	Error() error
}

type ErrorInjector interface {
	Error() error
}

type LatencyInjector interface {
	Latency()
}
