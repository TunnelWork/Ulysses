package logger

import (
	"fmt"
	"os"
)

// Non-block
func _Debug(v ...interface{}) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	go _debug(v...)
}

// Non-block
func _Info(v ...interface{}) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	go _info(v...)
}

// Non-block
func _Warning(v ...interface{}) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	go _warning(v...)
}

// Non-block
func _Error(v ...interface{}) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	go _error(v...)
}

// Block!
func _Fatal(v ...interface{}) {
	_fatal(v...) // Not calling as goroutine because non-block.
}

func _debug(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print("DEBUG: ", fmt.Sprint(v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("DEBUG: ", fmt.Sprint(v...), "\n")
	}
}

func _info(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print("INFO: ", fmt.Sprint(v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("INFO: ", fmt.Sprint(v...), "\n")
	}
}

func _warning(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print("WARNING: ", fmt.Sprint(v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("WARNING: ", fmt.Sprint(v...), "\n")
	}
}

func _error(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print("ERROR: ", fmt.Sprint(v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("ERROR: ", fmt.Sprint(v...), "\n")
	}
}

func _fatal(v ...interface{}) {
	loggerMutex.Lock()

	if verboseLogging {
		fmt.Print("FATAL: ", fmt.Sprint(v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("FATAL: ", fmt.Sprint(v...), "\n")
	}

	loggerMutex.Unlock()

	if exitingFunc != nil { // if set, call exitingFunc() first to clear goroutines
		(*exitingFunc)()
	}

	if loggerWaitGroup != nil {
		loggerWaitGroup.Wait() // Fatal will exit the system, so make sure the WaitGroup is cleared before that.
	}
	os.Exit(1)
}
