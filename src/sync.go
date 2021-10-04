package main

import (
	"sync"

	"github.com/TunnelWork/Ulysses/src/internal/logger"
)

var (
	// globalWaitGroup.Add(1) before starting a persistent goroutine.
	// globalWaitGroup.Done() before exiting a persistent goroutine.
	globalWaitGroup sync.WaitGroup

	// globalTickGroup.Wait() in every goroutine relying on system ticking:
	// - Any READ calls in server.go
	//
	globalTickGroup   sync.WaitGroup
	globalExitChannel chan bool
)

func initWaitGroup() {
	globalWaitGroup = sync.WaitGroup{}
	globalExitChannel = make(chan bool)
}

// globalExitSignal() should write to the globalExitChannel
// in order to quit all persistent goroutines, e.g., systemTicking
// it is very dirty now
func globalExitSignal() {
	// Signal up to 100 routines to quit.
	// globalWaitGroup.Add(1)
	logger.Debug("globalExitSignal(): everybody get out!")
	go func() {
		// defer globalWaitGroup.Done()
		for i := 0; i < 100; i++ {
			globalExitChannel <- true
		}
	}()
}
