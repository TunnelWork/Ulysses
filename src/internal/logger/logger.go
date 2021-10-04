package logger

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/TunnelWork/Ulysses/src/internal/conf"
)

var (
	verboseLogging  bool        = false
	loggerMutex     *sync.Mutex // TODO: Use a ticket lock for fairness, especially at high concurrency
	fileLogger      *log.Logger = nil
	loggerWaitGroup *sync.WaitGroup
	exitingFunc     *func()

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

func InitWithWaitGroupAndExitingFunc(wg *sync.WaitGroup, exitFunc *func(), loggerConfig conf.LoggerConfig) error {
	loggerWaitGroup = wg
	exitingFunc = exitFunc
	return Init(loggerConfig)
}

// LastWord() blocks and is only used for CLEAN, INTENDED EXITING.
// calling LastWord() does not invoke exitingFunc()
func LastWord(v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		loggerWaitGroup.Wait()
	}
	if verboseLogging {
		fmt.Print("LASTWORD: ", fmt.Sprint(v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("LASTWORD: ", fmt.Sprint(v...), "\n")
	}
	os.Exit(0)
}

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

	loggerWaitGroup.Wait() // Fatal will exit the system, so make sure the WaitGroup is cleared before that.
	os.Exit(1)
}
