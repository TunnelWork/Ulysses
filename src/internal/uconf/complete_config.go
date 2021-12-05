package uconf

import "github.com/TunnelWork/Ulysses.Lib/logging"

// Full Config defines all the configurable options for Ulysses.

// Should be loaded from MySQL Database

type CompleteConfig struct {
	// Mysql Client, should be internal?
	Mysql MysqlClientConfig // directly copy from file

	// Http Server
	Http HttpServerConfig // fetch from db

	// Logging
	Logger logging.LoggerConfig // fetch from db

	// Security
	Security SecurityModuleConfig // directly copy from file
}
