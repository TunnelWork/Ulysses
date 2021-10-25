package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/TunnelWork/Ulysses.Lib/server"
	"github.com/TunnelWork/Ulysses/src/internal/db"
	"github.com/TunnelWork/Ulysses/src/internal/logger"
	"github.com/gin-gonic/gin"
)

const (
	serverConfigTableName                           = `Servers`
	mysqlConnectTimeout                             = 500 * time.Millisecond // If not established within timeout, will fail.
	reloadUlyssesServerSignature tickEventSignature = 0xFEEDBEEF             // temp, we will come up with better names
)

// ServerManager is capable of performing MySQL CRUD operations to
// - Server Configuration Info
type ServerManager struct {
	dbConnector *db.MysqlConnector
}

func NewServerManager(dbconf db.DatabaseConfig) *ServerManager {
	return &ServerManager{
		dbConnector: db.NewMysqlConnector(dbconf),
	}
}

// Add() returns the newly inserted serverconf (id, nil), or (0, error) if any error
func (sm *ServerManager) Add(serverType string, confJson server.Configurables) (id uint, err error) {
	serverConfBytes, err := json.Marshal(confJson)
	if err != nil {
		// logger.Error("*ServerManager.Add(): cannot json.Marshal(), error: ", err)
		return id, err
	}

	dbConn, err := sm.dbConnector.Conn()

	if err != nil {
		// logger.Error("*ServerManager.Add(): cannot connect DB, error: ", err)
		return id, err
	}
	defer dbConn.Close()

	stmtInsertServerConf, err := dbConn.Prepare(`INSERT INTO ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` (ServerType, ConfJson, LastUpdate) VALUES( ?, ?, ? )`)
	if err != nil {
		// logger.Error("*ServerManager.Add(): cannot prepare statement. error: ", err)
		return id, err
	}
	defer stmtInsertServerConf.Close()

	result, err := stmtInsertServerConf.Exec(serverType, string(serverConfBytes), time.Now().Unix())
	if err != nil {
		// logger.Error("*ServerManager.Add(): cannot execute prepared statement. error: ", err)
		return id, err
	} else {
		idint64, err := result.LastInsertId()
		if err != nil {
			// logger.Error("*ServerManager.Add(): cannot get last inserted id. error: ", err)
			return id, err
		}
		id = uint(idint64)
	}

	logger.Info("*ServerManager.Add(): added new server with type ", serverType)
	return id, nil
}

func (sm *ServerManager) Lookup(id uint) (serverType string, confJson server.Configurables, err error) {
	confJson = server.Configurables{}

	dbConn, err := sm.dbConnector.Conn()
	if err != nil {
		// logger.Error("*ServerManager.Lookup(): cannot connect db, error: ", err)
		return "", server.Configurables{}, err
	}
	defer dbConn.Close()

	stmtLookupServerConf, err := dbConn.Prepare(`SELECT ServerType, ConfJson FROM ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` WHERE ID = ? AND Disabled = 0 AND DeletionTime = 0`)
	if err != nil {
		// logger.Error("*ServerManager.Lookup(): cannot prepare statement. error: ", err)
		return "", server.Configurables{}, err
	}

	var confJsonStr string

	err = stmtLookupServerConf.QueryRow(id).Scan(&serverType, &confJsonStr)
	if err != nil {
		// if err != sql.ErrNoRows { // Expect to see ErrNoRows a lot. Not even an error.
		// 	logger.Error("*ServerManager.Lookup(): can't query or scan. error: ", err)
		// } else {
		// 	logger.Error("*ServerManager.Lookup(): no such id or inactive: ", id)
		// }
		return "", server.Configurables{}, err
	}

	json.Unmarshal([]byte(confJsonStr), &confJson)

	return serverType, confJson, nil
}

func (sm *ServerManager) Update(id uint, serverType string, confJson server.Configurables) error {
	serverConfBytes, err := json.Marshal(confJson)
	if err != nil {
		logger.Error("*ServerManager.Update(): cannot json.Marshal(), error: ", err)
		return err
	}

	dbConn, err := sm.dbConnector.Conn()
	if err != nil {
		logger.Error("*ServerManager.Update(): cannot connect DB, error: ", err)
		return err
	}
	defer dbConn.Close()

	stmtUpdateServerConf, err := dbConn.Prepare(`UPDATE ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` SET ServerType = ?, ConfJson = ?, LastUpdate = ? WHERE ID = ? AND DeletionTime = 0`)
	if err != nil {
		logger.Error("*ServerManager.Add(): cannot prepare statement. error: ", err)
		return err
	}
	defer stmtUpdateServerConf.Close()

	_, err = stmtUpdateServerConf.Exec(serverType, string(serverConfBytes), time.Now().Unix(), id)
	if err != nil {
		logger.Error("*ServerManager.Update(): cannot execute prepared statement, error: ", err)
		return err
	}

	logger.Info("*ServerManager.Update(): updated server with id ", id, " and type ", serverType)
	return nil
}

func (sm *ServerManager) Delete(id uint) error {
	dbConn, err := sm.dbConnector.Conn()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	stmtDeleteServerConf, err := dbConn.Prepare(`UPDATE ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` SET LastUpdate = ?, DeletionTime = ? WHERE ID = ? AND DeletionTime = 0`)
	if err != nil {
		return err
	}
	defer stmtDeleteServerConf.Close()

	deletionTime := time.Now().Unix()
	_, err = stmtDeleteServerConf.Exec(deletionTime, deletionTime, id)
	if err != nil {
		return err
	}

	logger.Info("*ServerManager.Delete(): removed server with id ", id)
	return nil
}

func (sm *ServerManager) Disable(id uint) error {
	dbConn, err := sm.dbConnector.Conn()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	stmtDisableServerConf, err := dbConn.Prepare(`UPDATE ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` SET LastUpdate = ?, Disabled = 1 WHERE ID = ? AND Disabled = 0 AND DeletionTime = 0`)
	if err != nil {
		return err
	}
	defer stmtDisableServerConf.Close()

	_, err = stmtDisableServerConf.Exec(time.Now().Unix(), id)
	if err != nil {
		return err
	}

	logger.Info("*ServerManager.Disable(): disabled server with id ", id)
	return nil
}

func (sm *ServerManager) Enable(id uint) error {
	dbConn, err := sm.dbConnector.Conn()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	stmtDisableServerConf, err := dbConn.Prepare(`UPDATE ` + masterConfig.DB.TblPrefix + serverConfigTableName + ` SET LastUpdate = ?, Disabled = 0 WHERE ID = ? AND Disabled = 1 AND DeletionTime = 0`)
	if err != nil {
		return err
	}
	defer stmtDisableServerConf.Close()

	_, err = stmtDisableServerConf.Exec(time.Now().Unix(), id)
	if err != nil {
		return err
	}

	logger.Info("*ServerManager.Disable(): disabled server with id ", id)
	return nil
}

func _debugHandlerServerManager(c *gin.Context) {
	debugServerManager := NewServerManager(masterConfig.DB)
	op := c.Query("op")

	switch op {
	case "add":
		var confJson = server.Configurables{}
		err := json.Unmarshal([]byte(c.PostForm("conf_json")), &confJson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
		} else {
			id, err := debugServerManager.Add(c.PostForm("server_type"), confJson)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"op":     op,
					"status": "error",
					"error":  err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"op":     op,
					"status": "success",
					"id":     id,
				})
			}
		}
	case "lookup":
		var id uint
		iduint64, err := strconv.ParseUint(c.Query("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
		} else {
			id = uint(iduint64)
			serverType, confJson, err := debugServerManager.Lookup(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"op":     op,
					"status": "error",
					"error":  err.Error(),
				})
			} else {
				confJsonByte, err := json.Marshal(confJson)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"op":     op,
						"status": "error",
						"error":  err.Error(),
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"op":          op,
						"status":      "success",
						"server_type": serverType,
						"conf_json":   string(confJsonByte),
					})
				}
			}
		}
	case "update":
		var id uint
		iduint64, err := strconv.ParseUint(c.PostForm("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		var confJson = server.Configurables{}
		err = json.Unmarshal([]byte(c.PostForm("conf_json")), &confJson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		id = uint(iduint64)
		err = debugServerManager.Update(id, c.PostForm("server_type"), confJson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"op":     op,
			"status": "success",
		})
	case "delete":
		var id uint
		iduint64, err := strconv.ParseUint(c.PostForm("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
		}
		id = uint(iduint64)
		err = debugServerManager.Delete(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"op":     op,
			"status": "success",
		})
	case "disable":
		var id uint
		iduint64, err := strconv.ParseUint(c.PostForm("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
		}
		id = uint(iduint64)
		err = debugServerManager.Disable(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"op":     op,
			"status": "success",
		})
	case "enable":
		var id uint
		iduint64, err := strconv.ParseUint(c.PostForm("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
		}
		id = uint(iduint64)
		err = debugServerManager.Enable(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"op":     op,
				"status": "error",
				"error":  err.Error(),
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"op":     op,
			"status": "success",
		})
	default:

	}
}
