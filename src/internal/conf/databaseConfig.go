package conf

type DatabaseConfig struct {
	// Mandatory
	Host     string `yaml:"host"` // For IPv6, use the format of [::]
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Passwd   string `yaml:"passwd"`

	// Optional
	CA         string `yaml:"ca_cert"`     // Required by ClientKey & ClientCert
	ClientKey  string `yaml:"client_key"`  // Required by ClientCert
	ClientCert string `yaml:"client_cert"` // Required by ClientKey
}
