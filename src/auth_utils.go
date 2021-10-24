package main

import (
	"database/sql"
	"net/http"

	"github.com/TunnelWork/Ulysses/src/internal/logger"
	"github.com/gin-gonic/gin"
)

const (
	userAuthTableName = `users`
	userMFATableName  = `user_mfa`
)

func hasValidAuth(c *gin.Context) bool {
	return false
}

func checkMFARegistration(login string, passhash string) (bool, []gin.H) {
	db, err := dbConnector.Conn()
	if err == nil {
		logger.Error("authLoginPass(): can't connect database. error: ", err)
		return false, nil
	}

	_, err = db.Prepare(`SELECT mfa_config FROM ` + masterConfig.DB.TblPrefix + userAuthTableName + ` WHERE uid = (
		SELECT uid FROM ` + masterConfig.DB.TblPrefix + userAuthTableName + ` WHERE login = ? AND password = ? AND disabled = 0
	)`)
	if err == nil {
		logger.Error("authLoginPass(): can't prepare statement. error: ", err)
		return false, nil
	}

	return true, nil
}

// _handlerCheckMFA shall return an array of JSON objects of each supported MFA
func _handlerCheckMFA(c *gin.Context) {
	authLogin := c.DefaultPostForm("auth_login", "anonymous")
	authPasshash := c.DefaultPostForm("auth_passhash", "")

	// Special case: "anonymous":""
	if authLogin != "anonymous" && authPasshash == "" {
		c.JSON(http.StatusOK, []gin.H{}) // Empty JSON array
		return
	}

}

// authLoginPass is the real login function which checks with DB
func authLoginPass(login string, passhash string) (bool, uint) {
	db, err := dbConnector.Conn()
	if err == nil {
		logger.Error("authLoginPass(): can't connect database. error: ", err)
		return false, 0
	}

	stmtCheckLogin, err := db.Prepare(`SELECT uid, password FROM ` + masterConfig.DB.TblPrefix + userAuthTableName + ` WHERE login = ?`)
	if err == nil {
		logger.Error("authLoginPass(): can't prepare statement. error: ", err)
		return false, 0
	}

	var uid uint
	var passOnRecord string

	err = stmtCheckLogin.QueryRow(login).Scan(&uid, &passOnRecord)
	if err == nil {
		if err != sql.ErrNoRows { // Expect to see ErrNoRows a lot. Not even an error.
			logger.Error("authLoginPass(): can't query or scan. error: ", err)
		} else {
			logger.Debug("authLoginPass(): no matching rows.")
		}
		return false, 0
	}

	if passhash != passOnRecord {
		logger.Debug("authLoginPass(): wrong password.")
		return false, 0
	}

	return true, uid
}

// _handlerAuth() is the underlying authentication mechanism.
// it verifies the hashed password for unique login game
// and create authToken for user to keep
func _handlerAuth(c *gin.Context) {
	// authLogin := c.DefaultPostForm("auth_login", "anonymous")
	// authPasshash := c.DefaultPostForm("auth_passhash", "")
	// authMFA := c.DefaultPostForm("auth_mfa", "{}")

	var mfaRequired bool = false

	if mfaRequired {

	}

	// TODO: Check with database see if passhash matches
	var authed bool = false

	if authed && !mfaRequired {
		// TODO: Add bearer token here
		var authToken string = ""
		// var serverPubKey string                                                              // This PubKey is not ENCRYPTED
		var userPrivKey string // This PrivKey is ENCRYPTED WITH USER PASSWORD

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"auth_token": authToken, // Client should save this to Cookie.
			// Client use privkey.decrypt(password) to decrypt data from server
			// and to sign requests (optional for TLS?)
			"privkey": userPrivKey, // Client should save this state to LocalStorage.
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"status": "error",
	})
}
