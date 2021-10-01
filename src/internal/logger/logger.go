package logger

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/TunnelWork/Ulysses/src/internal/conf"
)

var (
	verboseLogging bool        = false
	loggerMutex    *sync.Mutex // TODO: Use a ticket lock for fairness, especially at high concurrency
	fileLogger     *log.Logger = nil

	Debug   = func(...interface{}) {} // Trivial and aligning with best practices
	Info    = func(...interface{}) {} // Non-trivial and aligning with best practices
	Warning = func(...interface{}) {} // Non-trivial and not aligning with best practices
	Error   = func(...interface{}) {} // Important and not in good condition, system can keep up
	Fatal   = func(...interface{}) {} // Important and not in good condition, system can't keep up
)

func Init(loggerConfig conf.LoggerConfig) error {
	switch loggerConfig.Level {
	case LvlDebug:
		Debug = _Debug
		fallthrough
	case LvlInfo:
		Info = _Info
		fallthrough
	case LvlWarning:
		Warning = _Warning
		fallthrough
	case LvlError:
		Error = _Error
		fallthrough
	case LvlFatal:
		Fatal = _Fatal
	default:
		return ErrBadLoggingLvl
	}
	verboseLogging = loggerConfig.Verbose
	if loggerConfig.Filepath != "" {
		f, err := os.OpenFile(loggerConfig.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err // ErrBadOpenFIle
		}
		fileLogger = log.New(f, "", log.LstdFlags)
	}
	loggerMutex = &sync.Mutex{}
	return nil
}

func _Debug(v ...interface{}) {
	go _debug(v...)
}

func _Info(v ...interface{}) {
	go _info(v...)
}

func _Warning(v ...interface{}) {
	go _warning(v...)
}

func _Error(v ...interface{}) {
	go _error(v...)
}

func _Fatal(v ...interface{}) {
	go _fatal(v...)
}

func _debug(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
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
	if verboseLogging {
		fmt.Print("ERROR: ", fmt.Sprint(v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("ERROR: ", fmt.Sprint(v...), "\n")
	}
}

func _fatal(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if verboseLogging {
		fmt.Print("FATAL: ", fmt.Sprint(v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("FATAL: ", fmt.Sprint(v...), "\n")
	}
	os.Exit(1)
}
