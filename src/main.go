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

	select {}

	// Tired? Let's get out of here
	// logger.Debug("globalExitSignal(): Now signal for exiting.")
	// globalExitSignal()
	// globalWaitGroup.Wait()

	// // Good Night
	// logger.LastWord("main(): Gute Nacht.") // Any last words? LOL
}

// bizLogic should start all business logics.
// e.g., Gin Webserver for APIs
func bizLogic() {
	go startGinRouter()
}
