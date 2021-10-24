package conf

const (
	defaultHost    string = "127.0.0.1"
	defaultPort    uint16 = 8080
	defaultUrlPath string = "/"
)

type SystemConfig struct {
	Host                        string `yaml:"api_host,omitempty"`
	Port                        uint16 `yaml:"api_port"`
	SystemTickPeriodMillisecond uint16 `yaml:"sys_tick_per_ms"` // 1~65535ms per system tick.
	UrlDomain                   string `yaml:"api_domain"`
	UrlPath                     string `yaml:"api_path"` // relative to api_domain. Should align with your WebUI setup.
}

func defaultSystemConfig() SystemConfig {
	return SystemConfig{
		Host:    defaultHost,
		Port:    defaultPort,
		UrlPath: defaultUrlPath,
	}
}
