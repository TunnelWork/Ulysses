package db

type DatabaseConfig struct {
	// Mandatory Connection Info
	Host     string `yaml:"host"` // For IPv6, use the format of [::]
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Passwd   string `yaml:"passwd"`

	// Optional Connection Info
	CA         string `yaml:"ca_cert"`     // Required by ClientKey & ClientCert
	ClientKey  string `yaml:"client_key"`  // Required by ClientCert
	ClientCert string `yaml:"client_cert"` // Required by ClientKey

	// Other MySQL Info
	MysqlAutoCommit bool `yaml:"auto_commit"`

	// Table Prefix
	TblPrefix string `yaml:"table_prefix"` // If unset, use default: "ulys_".
}
