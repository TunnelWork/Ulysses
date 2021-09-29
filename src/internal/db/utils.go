package db

import (
	"errors"

	"gopkg.in/yaml.v2"
)

var (
	ErrCannotAppendCert = errors.New("/src/internal/db: AppendCertsFromPEM() failed")
	ErrIncompleteConf   = errors.New("/src/internal/db: configuration isn't complete")
	ErrMySQLNoConn      = errors.New("/src/internal/db: cannot establish mysql connection")
)

type MysqlConf struct {
	// Mandatory
	mysqlHost     string `yaml:"host"` // For IPv6, use the format of [::]
	mysqlPort     uint16 `yaml:"port"`
	mysqlDatabase string `yaml:"database"`
	mysqlUser     string `yaml:"user"`
	mysqlPasswd   string `yaml:"passwd"`

	// Optional
	mysqlCAPath   string `yaml:"ca-cert-file"`
	mysqlKeyPath  string `yaml:"client-key-file"`
	mysqlCertPath string `yaml:"client-cert-file"`
}

func LoadMySqlConf(config []byte) MysqlConf {
	mconf := MysqlConf{}
	err := yaml.Unmarshal(config, &mconf)
	if err != nil {
		return MysqlConf{}
	}
	return mconf
}
