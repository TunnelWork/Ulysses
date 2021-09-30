package conf

type LoggerConfig struct {
	Verbose  bool   `yaml:"verbose"`
	Filepath string `yaml:"log_path"`
	Level    uint8  `yaml:"log_level"`
}

const (
	defaultVerbose      bool   = false
	defaultLogFilepath  string = "./Ulysses.log"
	defaultLoggingLevel uint8  = 3
)

func defaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Verbose:  defaultVerbose,      // Non-verbose by default
		Filepath: defaultLogFilepath,  // Default logging file
		Level:    defaultLoggingLevel, // Default logging level: WARNING and above
	}
}
