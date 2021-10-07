package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
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

func Init(loggerConfig LoggerConfig) error {
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

func InitWithWaitGroupAndExitingFunc(wg *sync.WaitGroup, exitFunc *func(), loggerConfig LoggerConfig) error {
	loggerWaitGroup = wg
	exitingFunc = exitFunc
	return Init(loggerConfig)
}

// LastWord() blocks and is only used for CLEAN, INTENDED EXITING.
// calling LastWord() does not invoke exitingFunc()
// caller should make sure the system is in a CLEAN state which
// exit(0) could be invoked without breaking any other systems.
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
