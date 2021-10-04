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
// 	PRIMARY KEY (id),
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
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	dbConn, err := db.DBConnectWithContext(ctx, masterConfig.DB)
	defer cancel()
	if err != nil {
		logger.Error("reloadUlyssesServer(): cannot connect db, error: ", err)
		return
	}
	defer dbConn.Close()

	stmtFetchServerConf, err := dbConn.Prepare(`SELECT id, server_type, conf_json FROM ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` WHERE deletion_time = 0`)
	if err != nil {
		logger.Error("reloadUlyssesServer(): cannot prepare statement. error: ", err)
		return
	}

	fetchedServerConfs, err := stmtFetchServerConf.Query()
	if err != nil {
		logger.Error("reloadUlyssesServer(): cannot query database. error: ", err)
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
			logger.Error("reloadUlyssesServer(): cannot scan result. error: ", err)
			continue
		}

		if _, ok := serverConfMap[key]; ok {
			logger.Error("reloadUlyssesServer(): repeated indexing key ", key)
			continue
		}

		err = json.Unmarshal([]byte(serverConfStr), &serverConf)
		if err != nil {
			logger.Error("reloadUlyssesServer(): cannot unmarshal conf_json: ", serverConfStr, ". error: ", err)
			continue
		}

		serverConfMap[key] = pairNameSconf{
			serverTypeName: serverTypeName,
			serverConf:     serverConf,
		}
		logger.Debug("reloadUlyssesServer(): loaded serverConfMap[", key, "] as a ", serverTypeName, " server")
	}
}

// addUlyssesServer() enforces thread-safety with the serverConfMapMutex.
// it DOES NOT update the local map!
// If no ticker setup, a consequent call to reloadUlyssesServer() is mandatory.
func addUlyssesServer(serverTypeName string, serverConf server.Configurables) {
	serverConfBytes, err := json.Marshal(serverConf)
	if err != nil {
		logger.Error("addUlyssesServer(): cannot json.Marshal(), error: ", err)
		return
	}
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	dbConn, err := db.DBConnectWithContext(ctx, masterConfig.DB)
	defer cancel()
	if err != nil {
		logger.Error("addUlyssesServer(): cannot connect DB, error: ", err)
		return
	}
	defer dbConn.Close()

	stmtInsertServerConf, err := dbConn.Prepare(`INSERT INTO ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` (server_type, conf_json) VALUES( ?, ? )`)
	if err != nil {
		logger.Error("addUlyssesServer(): cannot prepare statement. error: ", err)
		return
	}
	defer stmtInsertServerConf.Close()

	_, err = stmtInsertServerConf.Exec(serverTypeName, string(serverConfBytes))
	if err != nil {
		logger.Error("addUlyssesServer(): cannot execute prepared statement. error: ", err)
		return
	}

	logger.Info("addUlyssesServer(): added new server with type ", serverTypeName)
	serverConfMapDirty = true
}

// updateUlyssesServer() enforces thread-safety with the serverConfMapMutex.
// it DOES NOT update the local map!
// If no ticker setup, a consequent call to reloadUlyssesServer() is mandatory.
func updateUlyssesServer(id uint, serverTypeName string, serverConf server.Configurables) {
	serverConfBytes, err := json.Marshal(serverConf)
	if err != nil {
		logger.Error("updateUlyssesServer(): cannot json.Marshal(), error: ", err)
		return
	}
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	dbConn, err := db.DBConnectWithContext(ctx, masterConfig.DB)
	defer cancel()
	if err != nil {
		logger.Error("updateUlyssesServer(): cannot connect DB, error: ", err)
		return
	}
	defer dbConn.Close()

	stmtUpdateServerConf, err := dbConn.Prepare(`UPDATE` + masterConfig.DB.TblPrefix + serverConfigTableName + ` SET server_type = ?, conf_json = ? WHERE id = ?`)
	if err != nil {
		logger.Error("addUlyssesServer(): cannot prepare statement. error: ", err)
		return
	}
	defer stmtUpdateServerConf.Close()

	_, err = stmtUpdateServerConf.Exec(serverTypeName, string(serverConfBytes), id)
	if err != nil {
		logger.Error("updateUlyssesServer(): cannot execute prepared statement, error: ", err)
	}

	logger.Info("updateUlyssesServer(): updated server with id ", id, " and type ", serverTypeName)
	serverConfMapDirty = true
}

// removeUlyssesServer() enforces thread-safety with the serverConfMapMutex.
// it DOES NOT update the local map!
// If no ticker setup, a consequent call to reloadUlyssesServer() is mandatory.
func removeUlyssesServer(id uint, serverTypeName string) {
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	dbConn, err := db.DBConnectWithContext(ctx, masterConfig.DB)
	defer cancel()
	if err != nil {
		logger.Error("removeUlyssesServer(): cannot connect DB, error: ", err)
		return
	}
	defer dbConn.Close()

	stmtUpdateServerConf, err := dbConn.Prepare(`DELETE FROM ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` WHERE id = ? AND server_type = ?`)
	if err != nil {
		logger.Error("removeUlyssesServer(): cannot prepare statement. error: ", err)
		return
	}
	defer stmtUpdateServerConf.Close()

	_, err = stmtUpdateServerConf.Exec(id, serverTypeName)
	if err != nil {
		logger.Error("removeUlyssesServer(): cannot execute prepared statement, error: ", err)
	}

	logger.Info("removeUlyssesServer(): removed server with id ", id, " and type ", serverTypeName)
	serverConfMapDirty = true
}

// readUlyssesServerConfMap() enforces thread-safety with the serverConfMapMutex.
// it returns a copy of the map
func readUlyssesServerConfMap() map[uint](pairNameSconf) {
	globalTickGroup.Wait()
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()

	var copyMap = serverConfMap
	return copyMap
}

// searchUlyssesServerConf() enforces thread-safety with the serverConfMapMutex.
// it returns serverTypeName (string) and serverConf (server.Configurables)
func searchUlyssesServerConf(id uint) (string, server.Configurables) {
	globalTickGroup.Wait()
	serverConfMapMutex.Lock()
	defer serverConfMapMutex.Unlock()

	var serverTypeName string
	var sconf server.Configurables

	if mapEntry, ok := serverConfMap[id]; ok {
		serverTypeName = mapEntry.serverTypeName
		sconf = mapEntry.serverConf
	}
	return serverTypeName, sconf
}

func instantiateUlyssesServer(id uint) (string, server.Server) {
	var serverTypeName string
	var serverInstance server.Server

	if mapEntry, ok := serverConfMap[id]; ok {
		serverTypeName = mapEntry.serverTypeName
		tmpServerInstance, err := server.NewServerByType(serverTypeName, mapEntry.serverConf)
		if err == nil {
			serverInstance = tmpServerInstance
			logger.Debug("instantiateUlyssesServer(): successfully instantiated server id ", id)
		} else {
			logger.Error("instantiateUlyssesServer(): error when instantiate server id ", id, " type ", serverTypeName)
		}
	}

	return serverTypeName, serverInstance
}
