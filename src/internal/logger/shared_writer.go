package logger

import "fmt"

type levelWriter struct {
	level uint8
}

// Write() is blocked here.
func (lw *levelWriter) Write(p []byte) (n int, err error) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	switch lw.level {
	case LvlDebug:
		_debug(p)
	case LvlInfo:
		_info(p)
	case LvlWarning:
		_warning(p)
	case LvlError:
		_error(p)
	case LvlFatal:
		_fatal(p)
	default:
		return 0, ErrBadLoggingLvl
	}
	return len(p), nil
}

func DebugWriter() *levelWriter {
	return &levelWriter{
		level: LvlDebug,
	}
}

func InfoWriter() *levelWriter {
	return &levelWriter{
		level: LvlInfo,
	}
}

func WarningWriter() *levelWriter {
	return &levelWriter{
		level: LvlWarning,
	}
}

func ErrorWriter() *levelWriter {
	return &levelWriter{
		level: LvlError,
	}
}

func FatalWriter() *levelWriter {
	return &levelWriter{
		level: LvlFatal,
	}
}

type customWriter struct {
	prefix string
	suffix string
}

func NewCustomWriter(prefix string, suffix string) *customWriter {
	return &customWriter{
		prefix: prefix,
		suffix: suffix,
	}
}

// Write() is blocked here.
func (cw *customWriter) Write(p []byte) (n int, err error) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print(cw.prefix, " ", fmt.Sprint(string(p)), cw.suffix)
	}
	if fileLogger != nil {
		fileLogger.Print(cw.prefix, fmt.Sprint(string(p)), cw.suffix)
	}
	return len(p), nil
}
