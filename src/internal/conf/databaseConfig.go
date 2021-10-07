package conf

import "github.com/TunnelWork/Ulysses/src/internal/db"

const (
	defaultTblPrefix string = "ulys_"
)

func defaultDatabaseConfig() db.DatabaseConfig {
	return db.DatabaseConfig{
		TblPrefix: defaultTblPrefix,
	}
}
