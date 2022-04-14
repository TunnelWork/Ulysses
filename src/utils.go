package main

/******** Sync ********/

// slaveBlock() is called by goroutines to prevent main() from exiting before all goroutines are done
func slaveBlock() {
	masterWaitGroup.Add(1)
}

// slaveUnblock() is called by goroutines to notify main() that they are done
func slaveUnblock() {
	masterWaitGroup.Done()
}

// slaveWait() is called by goroutines to prevent new goroutines from running once exiting routine starts
func slaveWait() {
	slaveWaitGroup.Wait()
}

// masterBlock() is called by main() to block new goroutines from running once exiting routine starts
func masterBlock() {
	slaveWaitGroup.Add(1)
}

// masterUnblock() is called by main() to notify goroutines that they can start running
// TODO: This implementation is NOT USED in current implementation.
func masterUnblock() {
	slaveWaitGroup.Done()
}

// masterWait() is called by main() to wait for all goroutines to finish
func masterWait() {
	masterWaitGroup.Wait()
}
