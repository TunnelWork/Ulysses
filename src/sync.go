package main

import "sync"

var (
	globalWaitGroup   *sync.WaitGroup
	globalExitChannel chan bool
)

func initWaitGroup() {
	globalWaitGroup = &sync.WaitGroup{}
	globalExitChannel = make(chan bool)
}

// globalExitSignal() is very dirty now
func globalExitSignal() {
	// Signal up to 100 routines to quit.
	for i := 0; i < 100; i++ {
		globalExitChannel <- true
	}
}
