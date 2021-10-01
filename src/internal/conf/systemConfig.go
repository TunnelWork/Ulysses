package conf

const (
	defaultHost string = "127.0.0.1"
	defaultPort uint16 = 8080
)

type SystemConfig struct {
	Host string `yaml:"api_host,omitempty"`
	Port uint16 `yaml:"api_port"`
}

func defaultSystemConfig() SystemConfig {
	return SystemConfig{
		Host: defaultHost,
		Port: defaultPort,
	}
}
