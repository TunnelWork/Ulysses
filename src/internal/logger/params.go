package logger

import "errors"

const (
	LvlNull uint8 = iota
	LvlFatal
	LvlError
	LvlWarning
	LvlInfo
	LvlDebug
)

var (
	ErrBadLoggingLvl error = errors.New("internal/logger.Init(): invalid logLevel")
)
