package main

import (
	"encoding/json"

	"github.com/TunnelWork/Ulysses/src/internal/db"
	"github.com/TunnelWork/Ulysses/src/internal/logger"
	"github.com/TunnelWork/Ulysses/src/server"
)

/**********
 *	File: ulysses_server.go
 *	Author: Gaukas Wang <i@gaukas.wang>
 *
 *	Note:
 * 	This files define/declare all the functions needed to:
 *	- Load (from DB) a map of saved server config for Ulysses.Server as map[uint]Configurables
 *		- Add new config to that map[uint]Configurables
 *		- Remove config from that map[uint]Configurables
 *		- Update the db for corresponding changes
 *  - Instantiate a Configurable into Ulysses.Server
 *		- Call Ulysses.Server member functions
 */

const (
	serverConfigTableName = `servers`
)

// var serverConfigTableCreateQuery = `
// CREATE TABLE ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` IF NOT EXISTS (
// 	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
// 	deletion_time BIGINT NOT NULL DEFAULT 0,
// 	conf_json TEXT NOT NULL,
// 	upload BIGINT UNSIGNED NOT NULL DEFAULT 0,
// 	PRIMARY KEY (id),
// 	INDEX (password)
// )`

var (
	serverConfMap = map[uint]server.Configurables{} // Initialize as empty map
)

func loadUlyssesServerFromDB() {
	// DB Conn Livess Test
	dbConn, err := db.DBConnect(masterConfig.DB)
	if err != nil {
		logger.Error("loadUlyssesServerFromDB(): ", err)
		return
	}

	stmtFetchServerConf, err := dbConn.Prepare(`SELECT id, conf_json FROM ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` WHERE deletion_time = 0`)
	if err != nil {
		logger.Error("loadUlyssesServerFromDB: cannot prepare statement. error: ", err)
		return
	}

	fetchedServerConfs, err := stmtFetchServerConf.Query()
	if err != nil {
		logger.Error("loadUlyssesServerFromDB: cannot query database. error: ", err)
		return
	}

	for fetchedServerConfs.Next() {
		var key uint
		var strConf string
		var conf server.Configurables = server.Configurables{}
		err = fetchedServerConfs.Scan(&key, &strConf)
		if err != nil {
			logger.Error("loadUlyssesServerFromDB: cannot scan result. error: ", err)
			continue
		}

		err = json.Unmarshal([]byte(strConf), &conf)
		if err != nil {
			logger.Error("loadUlyssesServerFromDB: cannot unmarshal conf_json: ", strConf, ". error: ", err)
			continue
		}

		if _, ok := serverConfMap[key]; ok {
			logger.Error("loadUlyssesServerFromDB: repeated indexing key ", key)
			continue
		}

		serverConfMap[key] = conf
		logger.Debug("loadUlyssesServerFromDB: loaded serverConfMap[", key, "] = ", strConf)
	}
}
