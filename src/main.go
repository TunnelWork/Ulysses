package main

import (
	"github.com/TunnelWork/Ulysses/src/internal/logger"
)

func main() {
	// Initialize Business Logic Here.
	logger.Debug("In main()")
	bizLogic()
	logger.Debug("Set ticking...")
	startSystemTicking() // start system ticking so everything really starts

	// Block until...
	select {
	// Internal Exiting Signal, or...
	case <-globalExitChannel:
		logger.Warning("main(): received on globalExitChannel. Maybe globalExitSignal() is called? Executing shutting down procedure. ")
	// External Exiting Signal
	case <-sysSig:
		logger.Warning("main(): SIGINT/SIGTERM received. Executing shutting down procedure. ")
	}

	globalExitSignal() // globalExitChannel <- true
	globalWaitGroup.Wait()

	logger.LastWord("main(): Gute Nacht.") // Any last words? LOL
}

// bizLogic should start all business logics.
// e.g., Gin Webserver for APIs
func bizLogic() {
	go startGinRouter()
}
