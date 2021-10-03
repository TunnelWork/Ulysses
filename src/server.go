package main

/**********
 *	File: server.go
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

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/TunnelWork/Ulysses/src/internal/db"
	"github.com/TunnelWork/Ulysses/src/internal/logger"
	"github.com/TunnelWork/Ulysses/src/server"
)

// var serverConfigTableCreateQuery = `
// CREATE TABLE ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` IF NOT EXISTS (
// 	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
//	server_type VARCHAR(32) NOT NULL,
// 	deletion_time BIGINT NOT NULL DEFAULT 0,
// 	conf_json TEXT NOT NULL,
// 	upload BIGINT UNSIGNED NOT NULL DEFAULT 0,
// 	PRIMARY KEY (id),
// 	INDEX (password)
// )`

type pairNameSconf struct {
	serverTypeName string
	serverConf     server.Configurables
}

const (
	serverConfigTableName = `servers`

	reloadUlyssesServerSignature tickEventSignature = 0xFEEDBEEF // temp, we will come up with better names
)

var (
	serverConfMapMutex *sync.Mutex              // Want to use RWMutex, but goroutine's RLock() is exclusive as well
	serverConfMap      map[uint](pairNameSconf) // Initialize as empty map
	serverConfMapDirty bool
)

func initUlyssesServer() {
	if serverConfMapMutex == nil {
		serverConfMapMutex = &sync.Mutex{}
	}
	serverConfMapDirty = true
	reloadUlyssesServer()

	registerTickEvent(reloadUlyssesServerSignature, reloadUlyssesServer)
}

// reloadUlyssesServer() enforces thread-safety with the serverConfMapMutex.
func reloadUlyssesServer() {
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()
	if !serverConfMapDirty {
		logger.Debug("reloadUlyssesServer(): map isn't dirty, skipping")
		return
	}

	serverConfMap = map[uint](pairNameSconf){} // Must clear the map
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	dbConn, err := db.DBConnectWithContext(ctx, masterConfig.DB)
	cancel()
	if err != nil {
		logger.Error("reloadUlyssesServer(): ", err)
		return
	}

	stmtFetchServerConf, err := dbConn.Prepare(`SELECT id, name, conf_json FROM ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` WHERE deletion_time = 0`)
	if err != nil {
		logger.Error("reloadUlyssesServer: cannot prepare statement. error: ", err)
		return
	}

	fetchedServerConfs, err := stmtFetchServerConf.Query()
	if err != nil {
		logger.Error("reloadUlyssesServer: cannot query database. error: ", err)
		return
	}

	serverConfMapDirty = false

	for fetchedServerConfs.Next() {
		var key uint
		var serverTypeName string
		var serverConfStr string
		var serverConf server.Configurables = server.Configurables{}
		err = fetchedServerConfs.Scan(&key, &serverTypeName, &serverConfStr)
		if err != nil {
			logger.Error("reloadUlyssesServer: cannot scan result. error: ", err)
			continue
		}

		if _, ok := serverConfMap[key]; ok {
			logger.Error("reloadUlyssesServer: repeated indexing key ", key)
			continue
		}

		err = json.Unmarshal([]byte(serverConfStr), &serverConf)
		if err != nil {
			logger.Error("reloadUlyssesServer: cannot unmarshal conf_json: ", serverConfStr, ". error: ", err)
			continue
		}

		serverConfMap[key] = pairNameSconf{
			serverTypeName: serverTypeName,
			serverConf:     serverConf,
		}
		logger.Debug("reloadUlyssesServer: loaded serverConfMap[", key, "] as a ", serverTypeName, " server")
	}
}

// addUlyssesServer() enforces thread-safety with the serverConfMapMutex.
// it DOES NOT update the local map!
// If no ticker setup, a consequent call to reloadUlyssesServer() is mandatory.
func addUlyssesServer() {
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()
	serverConfMapDirty = true
}

// updateUlyssesServer() enforces thread-safety with the serverConfMapMutex.
// it DOES NOT update the local map!
// If no ticker setup, a consequent call to reloadUlyssesServer() is mandatory.
func updateUlyssesServer() {
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()
	serverConfMapDirty = true
}

// removeUlyssesServer() enforces thread-safety with the serverConfMapMutex.
// it DOES NOT update the local map!
// If no ticker setup, a consequent call to reloadUlyssesServer() is mandatory.
func removeUlyssesServer() {
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()
	serverConfMapDirty = true
}
