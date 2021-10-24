package conf

import "github.com/TunnelWork/Ulysses/src/internal/db"

const (
	defaultTblPrefix string = "ulys_"

	defaultMysqlAutoCommit bool = true
)

func defaultDatabaseConfig() db.DatabaseConfig {
	return db.DatabaseConfig{
		MysqlAutoCommit: defaultMysqlAutoCommit,
		TblPrefix:       defaultTblPrefix,
	}
}
