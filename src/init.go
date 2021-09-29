package main

import (
	"flag"

	"github.com/TunnelWork/Ulysses/src/internal/db"
)

// Logging Level
// 0: Absolute No Logging, you don't even know how it crashed (unless panic)
// 5: Log every non-trivial thing. FBI Open Up!
const (
	logNull uint8 = iota
	logFatal
	logError
	logWarning
	logInfo
	logDebug
)

var (
	// Configs
	ulyssesConfigPath string
	ulyssesConfig     config = config{
		logLevel: logDebug,
	}
	dbConfigPath string
	dbConfig     db.MysqlConf
)

func init() {
	//// GLOBAL INIT BEGIN ////
	flag.StringVar(&ulyssesConfigPath, "ulyssesConfigPath", "./conf/ulysses.yaml", "Ulysses General Configuration File")
	flag.StringVar(&dbConfigPath, "dbConfigPath", "./conf/mysql.yaml", "Database (MySQL/MariaDB) Connection Configuration File")
	flag.Parse()
	initConfig()
	// Load ulyssesConfig

	/*** GLOBAL INIT END ***/
	/*** SYSTEM MODULE INIT BEGIN ***/
	initLogger() // First thing first or we have nowhere to write
	initDB()
	/*** SYSTEM MODULE INIT END ***/
}
