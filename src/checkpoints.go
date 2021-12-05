package main

/**************************************
***   File: src\checkpoints.go
***   Summary: Defines a type of function named checkpoint which validates the HTTP Request received on a route.
***	  Author: Gaukas Wang
***************************************/

import (
	"net"
	"net/http"

	themis "github.com/TunnelWork/Themis"
	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/TunnelWork/Ulysses.Lib/logging"
	"github.com/TunnelWork/Ulysses/src/driver"
	"github.com/TunnelWork/Ulysses/src/internal/utils"
	"github.com/gin-gonic/gin"
)

/*********************************************************************
 * All checkpoints available are defined here.
 *
 * For a checkpoint:
 * - The name should be in CamelCase and should be a statement claiming a criteria being met.
 * - The checkpoint itself, as a function variable, should be Exported.
 * - Should be of type `checkpoint`
 *********************************************************************/

// A checkpoint is a handler checks if the request meets a specific criteria and:
// - Calls c.AbortWithStatusJSON() when FAIL.
// - Return silently when PASS.
type checkpoint = handler

// Authorization Token
var (
	// 	AuthorizationMustBeValid checks if the request comes with a valid Authorization header
	// 	The Authorization header must be in the form of "Bearer <token>" where <token> parses to a valid JWT token.
	//	Request method: GET/POST
	// 	Response:
	//		401 - "ACCESS_TOKEN_INVALID", "ACCESS_TOKEN_REQUIRED"
	//
	AuthorizationMustBeValid checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		// Get Authorization token
		authHeader := c.Request.Header["Authorization"]
		if len(authHeader) != 1 {
			logging.Debug("Unexpected Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespAccessTokenRequired)
			return
		}
		authToken := authHeader[0]

		bearer, err := themis.ImportBearerToken(authToken, utils.TokenRevoker)
		if err != nil {
			logging.Debug("Failed to import bearer token: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespAccessTokenInvalid)
			return
		}
		err = bearer.Verify(utils.TokenPubKey)
		if err != nil {
			logging.Debug("Failed to verify bearer token: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespAccessTokenInvalid)
			return
		}
		// Additional check: IP address. Preventing Authorization hijacking.
		ab := bearer.Body()
		if !ab.IpAddr.Equal(net.ParseIP(c.ClientIP())) {
			logging.Debug("IP mismatch: %s, %s", ab.IpAddr.String(), c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespAccessTokenInvalid)
			return
		}
	}
	AuthorizationMustBeValidIfExists checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		// Get Authorization token
		authHeader := c.Request.Header["Authorization"]
		if len(authHeader) == 0 {
			c.Next()
			return
		}
		AuthorizationMustBeValid(c)
	}
)

// MFA
var (
	// MFAMustBeEnabled enforces the user to HAVE at least 1 MFA method REGISTERED. If not, this method is set tp fail.
	//
	// Pre-requisite: AuthorizationMustBeValid
	//
	// Request method: GET/POST
	// Request post body: (empty)
	//
	// Response:
	//		401 - "ACCESS_TOKEN_INVALID"
	MFAMustBeEnabled checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		uid, err := utils.AuthorizationToUserID(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}

		_, err = auth.GetUserByID(uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}

		if !auth.AnyMFARegistered(uid) {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespMfaNotFound) // this user has no MFA enabled
			return
		}
		c.Next()
	}

	// 	MFARespMustBeValid requires the request to include a valid MFA response body
	// 	The POST body must be in the format below, where <response_body> is a MFA type-specific JSON object
	//  Client should not include the key "mfa" when submitting empty response.
	//
	// Pre-requisite: AuthorizationMustBeValid
	//
	// 	Request method: POST
	// 	Request post body:
	//	{
	//		...
	//		"mfa": {
	//		    "type": "...",
	//          "response": <response_body>,
	// 		}
	//		...
	// 	}
	//
	// Response:
	//		401 - "MFA_RESPONSE_INVALID"
	//		403 - "USER_NOT_FOUND", "MFA_RESPONSE_REQUIRED"
	//		500 - "INTERNAL_SERVER_ERROR"
	MFARespMustBeValid checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		uid, err := utils.AuthorizationToUserID(c)
		if err != nil {
			logging.Debug("Can't get userID: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		var form driver.FormMfaSubmitChallenge = driver.FormMfaSubmitChallenge{
			Mfa: &driver.MfaChallengeResponse{
				Type:     "",
				Response: map[string]string{},
			},
		}
		err = c.BindJSON(&form)
		if err != nil {
			if auth.AnyMFARegistered(uid) { // If user registered with MFA, must provide MFA response.
				logging.Debug("A valid MFA response is expected, but not received.")
				c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespMfaResponseRequired)
			} else {
				c.Next() // Not registered with any, therefore no MFA needs to nor should be provided.
			}
			return
		}

		err = auth.MFASubmitChallenge(form.Mfa.Type, uid, form.Mfa.Response)
		if err != nil {
			logging.Debug("Failed to submit challenge: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespMfaResponseInvalid)
			return
		}
		c.Next()
	}
)

// User Role
var (
	// UserMustBeAccountUser/UserMustBeProductUser/UserMustBeBillingUser checks if
	// the authed user is an affiliation account/product/billing user.
	//
	// Pre-requisite: AuthorizationMustBeValid
	//
	// Request method: GET/POST
	// Request post body: (empty)
	// Response:
	//		403 - "ACCESS_DENIED"
	//		500 - "INTERNAL_SERVER_ERROR"
	UserMustBeAccountUser checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		// check if Affiliation User
		if user.Role.Includes(auth.AFFILIATION_ACCOUNT_USER) {
			c.Next()
			return
		} else if user.Role.Includes(auth.GLOBAL_ADMIN) { // System Admin has all privileges
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeAccountAdmin: user(%d) isn't an affiliation account user nor a system admin.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}
	UserMustBeProductUser checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		// check if Affiliation User
		if user.Role.Includes(auth.AFFILIATION_PRODUCT_USER) {
			c.Next()
			return
		} else if user.Role.Includes(auth.GLOBAL_ADMIN) { // System Admin has all privileges
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeProductAdmin: user(%d) isn't an affiliation product user nor a system admin.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}
	UserMustBeBillingUser checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}
		// check if Affiliation User
		if user.Role.Includes(auth.AFFILIATION_BILLING_USER) {
			c.Next()
			return
		} else if user.Role.Includes(auth.GLOBAL_ADMIN) { // System Admin has all privileges
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeBillingAdmin: user(%d) isn't an affiliation billing user nor a system admin.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}

	// UserMustBeAccountAdmin/UserMustBeProductAdmin/UserMustBeBillingAdmin checks if
	// the authed user is an affiliation account/product/billing admin.
	//
	// Pre-requisite: AuthorizationMustBeValid
	//
	// Request method: GET/POST
	// Request post body: (empty)
	// Response:
	//		403 - "ACCESS_DENIED"
	//		500 - "INTERNAL_SERVER_ERROR"
	UserMustBeAccountAdmin checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		// check if Affiliation Admin
		if user.Role.Includes(auth.AFFILIATION_ACCOUNT_ADMIN) {
			c.Next()
			return
		} else if user.Role.Includes(auth.GLOBAL_ADMIN) { // System Admin has all privileges
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeAccountAdmin: user(%d) isn't an affiliation account admin nor a system admin.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}
	UserMustBeProductAdmin checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		// check if Affiliation Admin
		if user.Role.Includes(auth.AFFILIATION_PRODUCT_ADMIN) {
			c.Next()
			return
		} else if user.Role.Includes(auth.GLOBAL_ADMIN) { // System Admin has all privileges
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeProductAdmin: user(%d) isn't an affiliation product admin nor a system admin.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}
	UserMustBeBillingAdmin checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		// check if Affiliation Admin
		if user.Role.Includes(auth.AFFILIATION_BILLING_ADMIN) {
			c.Next()
			return
		} else if user.Role.Includes(auth.GLOBAL_ADMIN) { // System Admin has all privileges
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeBillingAdmin: user(%d) isn't an affiliation billing admin nor a system admin.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}

	// UserMustBeInEvaluation/UserMustBeInProduction checks if the user is in EVALUATION/PRODUCTION role.
	//
	// Pre-requisite: AuthorizationMustBeValid
	//
	// Request method: GET/POST
	// Request post body: (empty)
	UserMustBeInEvaluation checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		if user.Role.Includes(auth.GLOBAL_PRODUCTION_USER) {
			logging.Debug("UserMustBeInEvaluation: user(%d) is in PRODUCTION and therefore not considered as in EVALUATION.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		} else if user.Role.Includes(auth.GLOBAL_EVALUATION_USER) {
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeInEvaluation: user(%d) isn't in either EVALUATION nor PRODUCTION. MAYBE BROKEN USER?", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}
	UserMustBeInProduction checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		if user.Role.Includes(auth.GLOBAL_EVALUATION_USER) {
			logging.Debug("UserMustBeInProduction: user(%d) is in EVALUATION and therefore not considered as in PRODUCTION.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		} else if user.Role.Includes(auth.GLOBAL_PRODUCTION_USER) {
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeInProduction: user(%d) isn't in either EVALUATION nor PRODUCTION. MAYBE BROKEN USER?", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}

	// UserMustBeGlobalAdmin checks if the user is a global admin.
	//
	// Pre-requisite: AuthorizationMustBeValid
	//
	// Request method: GET/POST
	// Request post body: (empty)
	UserMustBeGlobalAdmin checkpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("UserMustBeProductAdmin: Can't get user(%d): %s", user.ID(), err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr) // other unhandled error
			return
		}

		if user.Role.Includes(auth.GLOBAL_ADMIN) {
			c.Next()
			return
		} else {
			logging.Debug("UserMustBeGlobalAdmin: user(%d) isn't a global admin.", user.ID())
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		}
	}
)
