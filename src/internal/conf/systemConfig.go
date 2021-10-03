package conf

const (
	defaultHost string = "127.0.0.1"
	defaultPort uint16 = 8080
)

type SystemConfig struct {
	Host                        string `yaml:"api_host,omitempty"`
	Port                        uint16 `yaml:"api_port"`
	SystemTickPeriodMillisecond uint16 `yaml:"sys_tick_per_ms"` // 1~65535ms per system tick.
}

func defaultSystemConfig() SystemConfig {
	return SystemConfig{
		Host: defaultHost,
		Port: defaultPort,
	}
}
