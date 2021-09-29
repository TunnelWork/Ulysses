package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	loggerMutex *sync.Mutex // TODO: Use a ticket lock for fairness, especially at high concurrency
	fileLogger  *log.Logger = nil

	Debug   = func(...interface{}) {}
	Info    = func(...interface{}) {}
	Warning = func(...interface{}) {} //
	Error   = func(...interface{}) {} // Error which is handled (or at least handlable)
	Fatal   = func(...interface{}) {} // Fatal Error which prevents the system from continue
)

// Must be called after GLOBAL INIT has been called
func initLogger() {
	switch ulyssesConfig.logLevel {
	case logDebug:
		Debug = _Debug
		fallthrough
	case logInfo:
		Info = _Info
		fallthrough
	case logWarning:
		Warning = _Warning
		fallthrough
	case logError:
		Error = _Error
		fallthrough
	case logFatal:
		Fatal = _Fatal
	}

	if ulyssesConfig.logFile != "" {
		f, err := os.OpenFile(ulyssesConfig.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fileLogger = log.New(f, "", log.LstdFlags)
	}
}

func _Debug(v ...interface{}) {
	go _debug(v)
}

func _Info(v ...interface{}) {
	go _info(v)
}

func _Warning(v ...interface{}) {
	go _warning(v)
}

func _Error(v ...interface{}) {
	go _error(v)
}

func _Fatal(v ...interface{}) {
	go _fatal(v)
}

func _debug(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if ulyssesConfig.verboseLogging {
		fmt.Print("DEBUG: ", v, "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("DEBUG: ", v, "\n")
	}
}

func _info(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if ulyssesConfig.verboseLogging {
		fmt.Print("INFO: ", v, "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("INFO: ", v, "\n")
	}
}

func _warning(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if ulyssesConfig.verboseLogging {
		fmt.Print("WARNING: ", v, "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("WARNING: ", v, "\n")
	}
}

func _error(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if ulyssesConfig.verboseLogging {
		fmt.Print("ERROR: ", v, "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("ERROR: ", v, "\n")
	}
}

func _fatal(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if ulyssesConfig.verboseLogging {
		fmt.Print("FATAL: ", v, "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("FATAL: ", v, "\n")
	}
	os.Exit(1)
}
