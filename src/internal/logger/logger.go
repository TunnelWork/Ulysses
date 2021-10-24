package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warning(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Writer(prefix string, suffix string) io.Writer
}

type dualLogger struct {
	config LoggerConfig // preserve for internal debugging purpose

	mutex     *sync.Mutex
	waitGroup *sync.WaitGroup

	fileLogger *log.Logger

	exitFunc *func()

	_debug   func(...interface{})
	_info    func(...interface{})
	_warning func(...interface{})
	_error   func(...interface{})
	_fatal   func(...interface{})
}

func DualLogger(conf LoggerConfig, wg *sync.WaitGroup, exFn *func()) *dualLogger {
	var dl dualLogger = dualLogger{
		config:    conf,
		waitGroup: wg,
		exitFunc:  exFn,
	}

	if loggerMutex != nil {
		dl.mutex = loggerMutex
	} else {
		dl.mutex = &sync.Mutex{}
	}

	if dl.config.Filepath != "" {
		f, err := os.OpenFile(dl.config.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil // ErrBadOpenFIle
		}
		dl.fileLogger = log.New(f, "", log.LstdFlags)
	}

	switch dl.config.Level {
	case LvlDebug:
		dl._debug = func(v ...interface{}) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print("Debug: ", fmt.Sprint(v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Debug: ", fmt.Sprint(v...), "\n")
			}
		}
		fallthrough
	case LvlInfo:
		dl._info = func(v ...interface{}) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print("Info: ", fmt.Sprint(v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Info: ", fmt.Sprint(v...), "\n")
			}
		}
		fallthrough
	case LvlWarning:
		dl._warning = func(v ...interface{}) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print("Warning: ", fmt.Sprint(v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Warning: ", fmt.Sprint(v...), "\n")
			}
		}
		fallthrough
	case LvlError:
		dl._error = func(v ...interface{}) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print("Error: ", fmt.Sprint(v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Error: ", fmt.Sprint(v...), "\n")
			}
		}
		fallthrough
	case LvlFatal:
		dl._fatal = func(v ...interface{}) {
			dl.mutex.Lock()

			if dl.config.Verbose {
				fmt.Print("Fatal: ", fmt.Sprint(v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Fatal: ", fmt.Sprint(v...), "\n")
			}

			dl.mutex.Unlock()

			if dl.exitFunc != nil {
				(*dl.exitFunc)()
			}

			if dl.waitGroup != nil {
				dl.waitGroup.Wait()
			}
			os.Exit(1)
		}
	default:
		return nil
	}

	return &dl
}

func (dl *dualLogger) Debug(v ...interface{}) {
	dl._debug(v)
}

func (dl *dualLogger) Info(v ...interface{}) {
	dl._info(v)
}

func (dl *dualLogger) Warning(v ...interface{}) {
	dl._warning(v)
}

func (dl *dualLogger) Error(v ...interface{}) {
	dl._error(v)
}

func (dl *dualLogger) Fatal(v ...interface{}) {
	dl._fatal(v)
}

func (dl *dualLogger) Writer(prefix string, suffix string) io.Writer {
	return &dualLoggerWriter{
		_write: func(p []byte) (n int, err error) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print(" ", fmt.Sprint(p), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print(prefix, fmt.Sprint(p), suffix)
			}
			return len(p), nil
		},
	}
}
