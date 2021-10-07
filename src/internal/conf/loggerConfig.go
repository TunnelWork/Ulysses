package conf

import "github.com/TunnelWork/Ulysses/src/internal/logger"

const (
	defaultVerbose      bool   = false
	defaultLogFilepath  string = "./Ulysses.log"
	defaultLoggingLevel uint8  = 3
)

func defaultLoggerConfig() logger.LoggerConfig {
	return logger.LoggerConfig{
		Verbose:  defaultVerbose,      // Non-verbose by default
		Filepath: defaultLogFilepath,  // Default logging file
		Level:    defaultLoggingLevel, // Default logging level: WARNING and above
	}
}
