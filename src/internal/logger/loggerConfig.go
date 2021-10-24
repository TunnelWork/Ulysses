package logger

type LoggerConfig struct {
	Verbose  bool   `yaml:"verbose"`
	Filepath string `yaml:"log_path"`
	Level    uint8  `yaml:"log_level"`
}
