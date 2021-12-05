package main

import (
	"fmt"

	"github.com/TunnelWork/Ulysses.Lib/logging"
)

func main() {
	// // Initialize Business Logic Here.
	// logger.Debug("In main()")
	bizLogic()

	// Block until...
	select {
	// External Exiting Signal
	case <-sysSig:
		logging.Warning("main(): SIGINT/SIGTERM received. Executing shutting down procedure. ")
	// Internal
	case <-globalExitChannel:
		logging.Warning("main(): readings on globalExitChannel. Executing shutting down procedure. ")
	}

	masterBlock() // Prevent new goroutines from starting
	masterWait()  // Wait until all goroutines are done

	logging.LastWord("main(): Gute Nacht.") // Any last words? LOL
}

// bizLogic should start all business logics NON-BLOCKING
// e.g., Gin Webserver for APIs
func bizLogic() {
	go ginRouter.Run(fmt.Sprintf("%s:%d", completeConfig.Http.HTTPHost, completeConfig.Http.HTTPPort))
	crontab.Start()
}
