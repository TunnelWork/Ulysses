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
	"github.com/gin-gonic/gin/binding"
)

var (
	Auth    = driver.Auth{}
	Billing = driver.Billing{}
	Server  = driver.Server{}
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
		err := c.ShouldBindBodyWith(&form, binding.JSON)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
			return
		}

		user, err := auth.GetUserByEmailPassword(form.Email, utils.HashPassword(form.Password))
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// Create Authorization Token
		token, err := themis.GetNewBearerToken(user.ID(), net.ParseIP(c.ClientIP()), time.Hour, utils.TokenRevoker)
		if err != nil {
			utils.HandleError(c, err)
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
			utils.HandleError(c, err)
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
			utils.HandleError(c, err)
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
			utils.HandleError(c, err)
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
			utils.HandleError(c, err)
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
			utils.HandleError(c, err)
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
			utils.HandleError(c, err)
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
	// BillingRecord
	// GET api/billing/billingrecord?cmd=<ListByWalletID|ListAll>, Authorization needed
	// POST api/billing/billingrecord?cmd=<Create>, Authorization needed
	GETBillingRecord endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "ListByWalletID":
			Billing.BillingRecord.ListByWalletID(c, user)
		case "ListAll":
			Billing.BillingRecord.ListAll(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTBillingRecord endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "Create":
			Billing.BillingRecord.Create(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}

	// ProductListingGroup
	// GET api/billing/productlistinggroup?cmd=<List|GetByID>
	// POST api/billing/productlistinggroup?cmd=<Create|Update|Delete>
	GETBillingProductListingGroup endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		cmd := c.Query("cmd")
		switch cmd {
		case "List":
			Billing.ProductListingGroup.List(c)
		case "GetByID":
			Billing.ProductListingGroup.GetByID(c)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTBillingProductListingGroup endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		cmd := c.Query("cmd")
		switch cmd {
		case "Create":
			Billing.ProductListingGroup.Create(c)
		case "Update":
			Billing.ProductListingGroup.Update(c)
		case "Delete":
			Billing.ProductListingGroup.Delete(c)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}

	// ProductListing
	// GET api/billing/productlisting?cmd=<ListByGroupID|GetAvailableByID|GetByID>, Authorization optional
	// POST api/billing/productlisting?cmd=<Create|Update|Delete|Hide|Unhide|Discontinue|Reactivate>
	GETBillingProductListing endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/
		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "ListByGroupID":
			Billing.ProductListing.ListByGroupID(c, user)
		case "GetAvailableByID":
			Billing.ProductListing.GetAvailableByID(c)
		case "GetByID":
			Billing.ProductListing.GetByID(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTBillingProductListing endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		cmd := c.Query("cmd")
		switch cmd {
		case "Create":
			Billing.ProductListing.Create(c)
		case "Update":
			Billing.ProductListing.Update(c)
		case "Delete":
			Billing.ProductListing.Delete(c)
		case "Hide":
			Billing.ProductListing.Hide(c)
		case "Unhide":
			Billing.ProductListing.Unhide(c)
		case "Discontinue":
			Billing.ProductListing.Discontinue(c)
		case "Reactivate":
			Billing.ProductListing.Reactivate(c)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}

	// Product
	// GET api/billing/product?cmd=<GetBySN|ListByID|ListByOwner|ListAll>, Authorization needed
	// POST api/billing/product?cmd=<CreateByListingID|Update|Terminate|ScheduleForTerminate>, Authorization needed
	GETBillingProduct endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "GetBySN":
			Billing.Product.GetBySN(c, user)
		case "ListByID":
			Billing.Product.ListByID(c, user)
		case "ListByOwner":
			Billing.Product.ListByOwner(c, user)
		case "ListAll":
			Billing.Product.ListAll(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTBillingProduct endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "CreateByListingID":
			Billing.Product.CreateByListingID(c, user)
		case "Update":
			Billing.Product.Update(c, user)
		case "Terminate":
			Billing.Product.Terminate(c, user)
		case "ScheduleForTerminate":
			Billing.Product.ScheduleForTerminate(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}

	// Wallet
	// GET api/billing/wallet?cmd=<View>, Authorization needed
	// POST api/billing/wallet?cmd=<Deposit|Withdraw|Enable|Disable>, Authorization needed
	GETBillingWallet endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "View":
			Billing.Wallet.View(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTBillingWallet endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "Deposit":
			Billing.Wallet.Deposit(c, user)
		case "Withdraw":
			Billing.Wallet.Withdraw(c, user)
		case "Enable":
			Billing.Wallet.Enable(c, user)
		case "Disable":
			Billing.Wallet.Disable(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
)

// Server
var (
	// ProvisioningServer
	//// ProvisioningAccount
	//// GET api/server/provisioning/account?cmd=GetBySN, Authorization needed
	//// POST api/server/provisioning/account?cmd=<Update|Delete|Suspend|Unsuspend|Refresh>, Authorization needed
	GETProvisioningAccount endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "GetBySN":
			Server.Provisioning.Account.GetBySN(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
	POSTProvisioningAccount endpoint = func(c *gin.Context) {
		/************ START GOROUTINE HEADER ************/
		slaveWait()
		slaveBlock()
		defer slaveUnblock()
		/************  END GOROUTINE HEADER  ************/

		user, err := utils.AuthorizationToUser(c)
		if err != nil {
			logging.Debug("Can't get user: %s", err.Error())
			utils.HandleError(c, err)
			return
		}

		cmd := c.Query("cmd")
		switch cmd {
		case "Update":
			Server.Provisioning.Account.Update(c, user)
		case "Delete":
			Server.Provisioning.Account.Delete(c, user)
		case "Suspend":
			Server.Provisioning.Account.Suspend(c, user)
		case "Unsuspend":
			Server.Provisioning.Account.Unsuspend(c, user)
		case "Refresh":
			Server.Provisioning.Account.Refresh(c, user)
		default:
			c.JSON(http.StatusBadRequest, utils.RespBadRequest)
		}
	}
)
