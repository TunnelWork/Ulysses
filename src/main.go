package main

import (
	"time"

	"github.com/TunnelWork/Ulysses/src/internal/logger"
)

func main() {
	// Initialize Business Logic Here.
	mockBizLogic()
	startSystemTicking() // start system ticking so everything starts

	// 10 second to make sure everything is working fine
	time.Sleep(10 * time.Second)

	// Tired? Let's get out of here
	globalExitSignal()

	// Block until all goroutine to return
	globalWaitGroup.Wait()

	// All done
	logger.Warning("main(): system exiting... good night.")
}

// BizLogic should NOT block.
func mockBizLogic() {

}
