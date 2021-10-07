package main

import (
	"time"

	"github.com/TunnelWork/Ulysses/src/internal/logger"
)

func main() {
	// Initialize Business Logic Here.
	bizLogic()
	startSystemTicking() // start system ticking so everything really starts

	// 10 second to make sure everything is working fine
	logger.Debug("time.Sleep(): Sleep for 6s.")
	time.Sleep(6 * time.Second)

	// TODO: REMOVE
	// logger.Fatal("OMG I crashed ;)")
	// select {}

	// Tired? Let's get out of here
	logger.Debug("globalExitSignal(): Now signal for exiting.")
	globalExitSignal()

	// All done
	logger.Warning("main(): system will exit upon globalWaitGroup becoming cleared.")

	globalWaitGroup.Wait()

	// Good Night
	logger.LastWord("main(): Gute Nacht.") // Any last words? LOL
}

// bizLogic should start all business logics.
// e.g., Gin Webserver for APIs
func bizLogic() {

}
