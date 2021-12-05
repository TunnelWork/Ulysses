package main

/**************************************
***   File: src\endpoints.go
***   Summary: API endpoints for FaaS (Function-as-a-Service).
***	  Author: Gaukas Wang
***************************************/

import (
	"net"
	"net/http"
	"time"

	themis "github.com/TunnelWork/Themis"
	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/TunnelWork/Ulysses.Lib/logging"
	"github.com/TunnelWork/Ulysses/src/driver"
	"github.com/TunnelWork/Ulysses/src/internal/utils"
	"github.com/gin-gonic/gin"
)

var (
	Auth = driver.Auth{}
)

// A endpoint is a handler assuming the HTTP request has already passed the authentication check (if needed) and:
// - Invoke a corresponding function.
// - Make a HTTP response with c.JSON().
type endpoint = handler

// Authorize
var (
	// Request a Authorization token for a specific user.
	// POST /api/authorize
	// {
	//   "email": <email>,
	//   "password": <raw_password>
	// }
	// Success Response Sample:
	// 	{
	//		status: "success",
	//		payload: "Bearer <A_BASE64_STRING>"
	//	}
	// TODO: Switch to appleboy/gin-jwt
	Authorize endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		var form FormAuthorize = FormAuthorize{}
		err := c.BindJSON(&form)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
			return
		}

		user, err := auth.GetUserByEmailPassword(form.Email, utils.HashPassword(form.Password))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}

		// Create Authorization Token
		token, err := themis.GetNewBearerToken(user.ID(), net.ParseIP(c.ClientIP()), time.Hour, utils.TokenRevoker)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}
		token.Sign(utils.TokenPrivKey)
		token.SetFullToken()
		c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, token.GetFullToken()))
	}
)

// Auth
var (
	// Affiliation
	// GET api/auth/affiliation?cmd=<cmd>[&id=<id>]
	// POST api/auth/affiliation?cmd=<cmd>
	GETAuthAffiliation endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespNotAuthorized)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "GetByID":
			Auth.Affiliation.GetByID(c, user)
		case "ParentAffiliation":
			Auth.Affiliation.ParentAffiliation(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTAuthAffiliation endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespNotAuthorized)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "Create":
			Auth.Affiliation.Create(c, user)
		case "Update":
			Auth.Affiliation.Update(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}

	// MFA
	// GET api/auth/mfa?cmd=Registered&type=<type>
	// POST api/auth/mfa?cmd=<InitSignUp|CompleteSignUP|NewChallenge|SubmitChallenge|Remove>
	GETAuthMFA endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		uid, err := utils.AuthorizationToUserID(c)
		if err != nil {
			logging.Debug("Can't get uid: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "Registered":
			Auth.MFA.Registered(c, uid)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTAuthMFA endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		uid, err := utils.AuthorizationToUserID(c)
		if err != nil {
			logging.Debug("Can't get uid: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "InitSignUp":
			Auth.MFA.InitSignUp(c, uid)
		case "CompleteSignUp":
			Auth.MFA.CompleteSignUp(c, uid)
		case "NewChallenge":
			Auth.MFA.NewChallenge(c, uid)
		case "SubmitChallenge":
			Auth.MFA.SubmitChallenge(c, uid)
		case "Remove":
			Auth.MFA.Remove(c, uid)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}

	// User
	// GET api/auth/user?cmd=<GetByID|ListByAffiliation|EmailExists|Info>[&id=<id>][&affiliation=<affiliation>][email=<email>]
	// POST api/auth/user?cmd=<Create|Update|Wipe|CreateInfo|UpdateInfo>
	GETAuthUser endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "GetByID":
			Auth.User.GetByID(c, user)
		case "List":
			Auth.User.List(c, user)
		case "ListByAffiliation":
			Auth.User.ListByAffiliation(c, user)
		case "EmailExists":
			Auth.User.EmailExists(c)
		case "Info":
			Auth.User.Info(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTAuthUser endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInternalErr)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "Create":
			Auth.User.Create(c, user) // user might be nil here
		case "Update":
			Auth.User.Update(c, user)
		case "Wipe":
			Auth.User.Wipe(c, user)
		case "CreateInfo":
			Auth.User.CreateInfo(c, user)
		case "UpdateInfo":
			Auth.User.UpdateInfo(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
)

// Billing
var (
// BillingProductListingGroup
)
