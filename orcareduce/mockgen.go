package orcareduce

//go:generate mockgen -destination=./mock/mock.go -package=mock github.com/mkuchenbecker/orcareduce/orcareduce Reactor,Logger,Handler,Precondition,ID,Director,ActorDataSink,Injector,ErrorInjector,LatencyInjector
