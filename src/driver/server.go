package driver

import (
	"net/http"
	"strconv"

	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/TunnelWork/Ulysses.Lib/billing"
	"github.com/TunnelWork/Ulysses.Lib/server"
	"github.com/TunnelWork/Ulysses/src/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	Provisioning ProvisioningServer
}

type ProvisioningServer struct {
	Account ProvisioningAccount
}

func productToProvisioningServer(product *billing.Product) (server.ProvisioningServer, error) {
	// Get ProductListing
	productListing, err := billing.SudoGetProductListingByID(product.ProductID)
	if err != nil {
		return nil, err
	}

	// Get Server Instance
	return server.NewProvisioningServer(productListing.ServerType, productListing.ServerInstanceID, productListing.ServerConfiguration)
}

type ProvisioningAccount struct{}

func (ProvisioningAccount) GetBySN(c *gin.Context, user *auth.User) {
	sn, err := strconv.ParseUint(c.Query("sn"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	product, err := billing.GetProductBySerialNumber(sn)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if !productIsAvailableToUser(product, user) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
	}

	serverInstance, err := productToProvisioningServer(product)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Get Account Info
	account, err := serverInstance.GetAccount(sn)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	credentials, err := account.Credentials()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	customerCredentials := credentials.Customer()
	adminCredentials := credentials.Admin()

	resources, err := account.Resources()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Compose the response
	resp := gin.H{
		"SerialNumber": sn,
		"Account": gin.H{
			"Credentials": gin.H{
				"Customer": customerCredentials,
			},
			"Resources": resources,
		},
	}
	if user.Role.Includes(auth.GLOBAL_ADMIN) || user.Role.Includes(auth.AFFILIATION_PRODUCT_ADMIN) {
		resp["Account"].(gin.H)["Credentials"].(gin.H)["Admin"] = adminCredentials
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, resp))
}

func (ProvisioningAccount) Update(c *gin.Context, user *auth.User) {
	var form FormUpdateProvisioningAccount = FormUpdateProvisioningAccount{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	product, err := billing.GetProductBySerialNumber(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if !productIsManagedByUser(product, user) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// Get Server Instance
	serverInstance, err := productToProvisioningServer(product)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	err = serverInstance.UpdateAccount(form.SerialNumber, form.AccountConfiguration)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}

func (ProvisioningAccount) Delete(c *gin.Context, user *auth.User) {
	var form FormDeleteProvisioningAccount = FormDeleteProvisioningAccount{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	product, err := billing.GetProductBySerialNumber(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if !productIsManagedByUser(product, user) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// Get Server Instance
	serverInstance, err := productToProvisioningServer(product)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	err = serverInstance.DeleteAccount(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}

func (ProvisioningAccount) Suspend(c *gin.Context, user *auth.User) {
	var form FormDeleteProvisioningAccount = FormDeleteProvisioningAccount{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	product, err := billing.GetProductBySerialNumber(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if !productIsManagedByUser(product, user) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// Get Server Instance
	serverInstance, err := productToProvisioningServer(product)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	err = serverInstance.SuspendAccount(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}

func (ProvisioningAccount) Unsuspend(c *gin.Context, user *auth.User) {
	var form FormDeleteProvisioningAccount = FormDeleteProvisioningAccount{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	product, err := billing.GetProductBySerialNumber(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if !productIsManagedByUser(product, user) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// Get Server Instance
	serverInstance, err := productToProvisioningServer(product)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	err = serverInstance.UnsuspendAccount(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}

func (ProvisioningAccount) Refresh(c *gin.Context, user *auth.User) {
	var form FormDeleteProvisioningAccount = FormDeleteProvisioningAccount{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	product, err := billing.GetProductBySerialNumber(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Only GLOBAL_ADMIN can do this
	if !user.Role.Includes(auth.GLOBAL_ADMIN) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// Get Server Instance
	serverInstance, err := productToProvisioningServer(product)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	err = serverInstance.RefreshAccount(form.SerialNumber)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}
