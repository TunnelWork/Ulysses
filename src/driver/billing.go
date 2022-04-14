package driver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/TunnelWork/Ulysses.Lib/billing"
	"github.com/TunnelWork/Ulysses.Lib/payment"
	"github.com/TunnelWork/Ulysses.Lib/server"
	"github.com/TunnelWork/Ulysses/src/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Billing struct {
	BillingRecord       BillingRecord
	ProductListingGroup ProductListingGroup
	ProductListing      ProductListing
	Product             Product
	Wallet              Wallet
}

type BillingRecord struct {
}

func (BillingRecord) Create(c *gin.Context, user *auth.User) {
	// Check if user is GLOBAL_ADMIN
	if !user.Role.Includes(auth.GLOBAL_ADMIN) {
		c.JSON(http.StatusForbidden, utils.RespAccessDenied)
	} else {
		c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
	}
}

func (BillingRecord) ListByWalletID(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}

func (BillingRecord) ListAll(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}

type ProductListingGroup struct {
}

func (ProductListingGroup) List(c *gin.Context) {
	ids, err := billing.ListProductListingGroupIDs() // TODO: hide hidden ones
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, ids))
}

func (ProductListingGroup) GetByID(c *gin.Context) {
	idstr := c.Query("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	productListingGroup, err := billing.GetProductListingGroupByID(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, productListingGroup))
}

func (ProductListingGroup) Create(c *gin.Context) {
	var form FormCreateProductListingGroup = FormCreateProductListingGroup{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	var productListingGroup billing.ProductListingGroup = billing.ProductListingGroup{
		ProductGroupName:        form.Name,
		ProductGroupDescription: form.Description,
		Hidden:                  form.Hidden,
	}
	if plgId, err := billing.NewProductListingGroup(productListingGroup); err != nil {
		utils.HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, plgId))
	}
}

func (ProductListingGroup) Update(c *gin.Context) {
	var form FormUpdateProductListingGroup = FormUpdateProductListingGroup{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	var productListingGroup billing.ProductListingGroup = billing.ProductListingGroup{
		ProductGroupID:          form.ID,
		ProductGroupName:        form.Name,
		ProductGroupDescription: form.Description,
		Hidden:                  form.Hidden,
	}
	if err := productListingGroup.Save(); err != nil {
		utils.HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, utils.RespOK)
	}
}

func (ProductListingGroup) Delete(c *gin.Context) {
	var form FormDeleteProductListingGroup = FormDeleteProductListingGroup{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if err := billing.DeleteProductListingGroupByID(form.ID); err != nil {
		utils.HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, utils.RespOK)
	}
}

type ProductListing struct {
}

func (ProductListing) ListByGroupID(c *gin.Context, user *auth.User) {
	// ListProductListingsByGroupID & SudoListProductListingsByGroupID
	groupIdStr := c.Query("group_id")
	groupId, err := strconv.ParseUint(groupIdStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	var productListings []*billing.ProductListing

	if user == nil || !user.Role.Includes(auth.GLOBAL_ADMIN) {
		// Everyone can see non-hidden, non-discontinued products
		productListings, err = billing.ListProductListingsByGroupID(groupId)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
	} else {
		// Global admins can see all products
		productListings, err = billing.SudoListProductListingsByGroupID(groupId)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
	}

	// re-populate a list of map[string]interface{} with exported fields & funcs
	var list []map[string]interface{}
	for _, productListing := range productListings {
		exported := map[string]interface{}{
			"ProductID":           productListing.ProductID(),
			"ProductGroupID":      productListing.ProductGroupID,
			"ProductName":         productListing.ProductName,
			"ProductDescription":  productListing.ProductDescription,
			"ServerType":          productListing.ServerType,
			"BillingOptions":      productListing.BillingOptions,
			"UsageBillingFactors": productListing.UsageBillingFactors,
			"Hidden":              productListing.Hidden,
			"Discontinued":        productListing.Discontinued,
		}
		list = append(list, exported)
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, list))
}
func (ProductListing) GetAvailableByID(c *gin.Context) {
	// GetProductListingByID, for everyone
	idstr := c.Query("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	productListing, err := billing.GetProductListingByID(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Export useful fields
	exported := map[string]interface{}{
		"ProductID":           productListing.ProductID(),
		"ProductGroupID":      productListing.ProductGroupID,
		"ProductName":         productListing.ProductName,
		"ProductDescription":  productListing.ProductDescription,
		"ServerType":          productListing.ServerType,
		"BillingOptions":      productListing.BillingOptions,
		"UsageBillingFactors": productListing.UsageBillingFactors,
		"Hidden":              productListing.Hidden,
		"Discontinued":        productListing.Discontinued,
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, exported))
}
func (ProductListing) GetByID(c *gin.Context, user *auth.User) {
	// SudoGetProductListingByID, for logged in users
	idstr := c.Query("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespAccessDenied)
		return
	}

	productListing, err := billing.SudoGetProductListingByID(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Export useful fields
	exported := map[string]interface{}{
		"ProductID":           productListing.ProductID(),
		"ProductGroupID":      productListing.ProductGroupID,
		"ProductName":         productListing.ProductName,
		"ProductDescription":  productListing.ProductDescription,
		"ServerType":          productListing.ServerType,
		"BillingOptions":      productListing.BillingOptions,
		"UsageBillingFactors": productListing.UsageBillingFactors,
		"Hidden":              productListing.Hidden,
		"Discontinued":        productListing.Discontinued,
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, exported))
}

func (ProductListing) Create(c *gin.Context) {
	var form FormCreateProductListing = FormCreateProductListing{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	var productListing *billing.ProductListing = &form.ProductListing
	if plId, err := productListing.Add(); err != nil {
		utils.HandleError(c, err)
		return
	} else {
		c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, plId))
	}
}
func (ProductListing) Update(c *gin.Context) {
	var form FormUpdateProductListing = FormUpdateProductListing{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	onRecord, err := billing.SudoGetProductListingByID(form.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Update the on-record product listing
	onRecord.ProductGroupID = form.ProductListing.ProductGroupID
	onRecord.ProductName = form.ProductListing.ProductName
	onRecord.ProductDescription = form.ProductListing.ProductDescription
	onRecord.ServerType = form.ProductListing.ServerType
	onRecord.ServerInstanceID = form.ProductListing.ServerInstanceID
	onRecord.ServerConfiguration = form.ProductListing.ServerConfiguration
	onRecord.BillingOptions = form.ProductListing.BillingOptions
	onRecord.UsageBillingFactors = form.ProductListing.UsageBillingFactors

	if err := onRecord.Save(); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}
func (ProductListing) Delete(c *gin.Context) {
	var form FormDeleteProductListing = FormDeleteProductListing{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if err := billing.DeleteProductListingByID(form.ID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}
func (ProductListing) Hide(c *gin.Context) {
	// Hide()
	var form FormToggleProductListing = FormToggleProductListing{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	onRecord, err := billing.SudoGetProductListingByID(form.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err = onRecord.Hide(); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}
func (ProductListing) Unhide(c *gin.Context) {
	// Unhide()
	var form FormToggleProductListing = FormToggleProductListing{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	onRecord, err := billing.SudoGetProductListingByID(form.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err = onRecord.Unhide(); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}
func (ProductListing) Discontinue(c *gin.Context) {
	// Discontinue()
	var form FormToggleProductListing = FormToggleProductListing{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	onRecord, err := billing.SudoGetProductListingByID(form.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err = onRecord.Discontinue(); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}
func (ProductListing) Reactivate(c *gin.Context) {
	// Reactivate()
	var form FormToggleProductListing = FormToggleProductListing{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	onRecord, err := billing.SudoGetProductListingByID(form.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err = onRecord.Reactivate(); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.RespOK)
}

// func (ProductListing) ListAvailable(c *gin.Context) {
// }

type Product struct{}

func productIsAvailableToUser(product *billing.Product, user *auth.User) bool {
	// Check if any of the criteria is met:
	// - the product belongs to current user
	// - the product belongs to the affiliation of current user, who is an AFFILIATION_PRODUCT_USER
	// - the user is a GLOBAL_ADMIN
	if product.OwnerUserID == user.ID() {
		return true
	} else if user.AffiliationID != 0 && user.AffiliationID == product.OwnerAffiliationID && user.Role.Includes(auth.AFFILIATION_PRODUCT_USER) {
		return true
	} else if user.Role.Includes(auth.GLOBAL_ADMIN) {
		return true
	}

	return false
}
func productIsManagedByUser(product *billing.Product, user *auth.User) bool {
	// Check if any of the criteria is met:
	// - the product belongs to current user
	// - the product belongs to the affiliation of current user, who is an AFFILIATION_PRODUCT_ADMIN
	// - the user is a GLOBAL_ADMIN
	if product.OwnerUserID == user.ID() {
		return true
	} else if user.AffiliationID != 0 && user.AffiliationID == product.OwnerAffiliationID && user.Role.Includes(auth.AFFILIATION_PRODUCT_ADMIN) {
		return true
	} else if user.Role.Includes(auth.GLOBAL_ADMIN) {
		return true
	}

	return false
}

func (Product) GetBySN(c *gin.Context, user *auth.User) {
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

	// Check if any of the criteria is met:
	// - the product belongs to current user
	// - the product belongs to the affiliation of current user, who is an AFFILIATION_PRODUCT_USER
	// - the user is a GLOBAL_ADMIN
	if !productIsAvailableToUser(product, user) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// Extract owner info
	var owner gin.H
	if product.OwnerUserID != 0 {
		owner = gin.H{
			"Type": "User",
			"ID":   product.OwnerUserID,
		}
	} else if product.OwnerAffiliationID != 0 {
		owner = gin.H{
			"Type": "Affiliation",
			"ID":   product.OwnerAffiliationID,
		}
	} else {
		owner = gin.H{
			"Type": "System",
		}
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, gin.H{
		"SerialNumber": product.SerialNumber(),
		"Owner":        owner,
		"ProductInfo": gin.H{
			"ID":              product.ProductID,
			"DateCreation":    (product.DateCreation()).Format("2006-01-02"),
			"DateLastBill":    (product.DateLastBill()).Format("2006-01-02"),
			"DateTermination": (product.DateTermination()).Format("2006-01-02"),
			"Terminated":      product.Terminated(),
			"WalletID":        product.WalletID,
			"BillingOption":   product.BillingOption,
		},
	}))
}
func (Product) ListByID(c *gin.Context, user *auth.User) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	// User must be GLOBAL_ADMIN
	if user == nil || !user.Role.Includes(auth.GLOBAL_ADMIN) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// List products
	products, err := billing.ListProductsByProductID(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// For every product, build exported data
	var exported []gin.H
	for _, product := range products {
		// Extract owner info
		var owner gin.H
		if product.OwnerUserID != 0 {
			owner = gin.H{
				"Type": "User",
				"ID":   product.OwnerUserID,
			}
		} else if product.OwnerAffiliationID != 0 {
			owner = gin.H{
				"Type": "Affiliation",
				"ID":   product.OwnerAffiliationID,
			}
		} else {
			owner = gin.H{
				"Type": "System",
			}
		}

		exported = append(exported, gin.H{
			"SerialNumber": product.SerialNumber(),
			"Owner":        owner,
			"ProductInfo": gin.H{
				"ID":              product.ProductID,
				"DateCreation":    (product.DateCreation()).Format("2006-01-02"),
				"DateLastBill":    (product.DateLastBill()).Format("2006-01-02"),
				"DateTermination": (product.DateTermination()).Format("2006-01-02"),
				"Terminated":      product.Terminated(),
				"WalletID":        product.WalletID,
				"BillingOption":   product.BillingOption,
			},
		})
	}
	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, exported))
}
func (Product) ListByOwner(c *gin.Context, currentUser *auth.User) {
	uidStr := c.DefaultQuery("uid", "0")
	aidStr := c.DefaultQuery("aid", "0")
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}
	aid, err := strconv.ParseUint(aidStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if currentUser == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	if uid != 0 {
		// For a PERSONAL_PRODUCT, the current user must be one of the following:
		// - the user of uid specified
		// - AFFILIATION_PRODUCT_ADMIN of the affiliation of which the user of uid is a member
		// - GLOBAL_ADMIN
		user, err := auth.GetUserByID(uid)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		if user.ID() != currentUser.ID() && (currentUser.AffiliationID == 0 || currentUser.AffiliationID != user.AffiliationID || !currentUser.Role.Includes(auth.AFFILIATION_PRODUCT_ADMIN)) && !currentUser.Role.Includes(auth.GLOBAL_ADMIN) {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		}

		// List products
		products, err := billing.ListUserProducts(uid)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// For every product, build exported data
		var exported []gin.H
		for _, product := range products {
			exported = append(exported, gin.H{
				"SerialNumber": product.SerialNumber(),
				"Owner": gin.H{
					"Type": "User",
					"ID":   uid,
				},
				"ProductInfo": gin.H{
					"ID":              product.ProductID,
					"DateCreation":    (product.DateCreation()).Format("2006-01-02"),
					"DateLastBill":    (product.DateLastBill()).Format("2006-01-02"),
					"DateTermination": (product.DateTermination()).Format("2006-01-02"),
					"Terminated":      product.Terminated(),
					"WalletID":        product.WalletID,
					"BillingOption":   product.BillingOption,
				},
			})
		}

		c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, exported))
	} else if aid != 0 {
		// For an AFFILIATION_PRODUCT, the current user must be one of the following:
		// - AFFILIATION_PRODUCT_USER of the target affiliation
		// - GLOBAL_ADMIN
		if (currentUser.AffiliationID != aid || !currentUser.Role.Includes(auth.AFFILIATION_PRODUCT_USER)) && !currentUser.Role.Includes(auth.GLOBAL_ADMIN) {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		}

		// List products
		products, err := billing.ListAffiliationProducts(aid)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// For every product, build exported data
		var exported []gin.H
		for _, product := range products {
			exported = append(exported, gin.H{
				"SerialNumber": product.SerialNumber(),
				"Owner": gin.H{
					"Type": "Affiliation",
					"ID":   aid,
				},
				"ProductInfo": gin.H{
					"ID":              product.ProductID,
					"DateCreation":    (product.DateCreation()).Format("2006-01-02"),
					"DateLastBill":    (product.DateLastBill()).Format("2006-01-02"),
					"DateTermination": (product.DateTermination()).Format("2006-01-02"),
					"Terminated":      product.Terminated(),
					"WalletID":        product.WalletID,
					"BillingOption":   product.BillingOption,
				},
			})
		}

		c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, exported))
	} else {
		// For a system owned product, the current user must be GLOBAL_ADMIN or GLOBAL_INTERNAL_USER
		if !currentUser.Role.Includes(auth.GLOBAL_ADMIN) && !currentUser.Role.Includes(auth.GLOBAL_INTERNAL_USER) {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		}

		// List products
		products, err := billing.ListSystemProducts()
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// For every product, build exported data
		var exported []gin.H
		for _, product := range products {
			exported = append(exported, gin.H{
				"SerialNumber": product.SerialNumber(),
				"Owner": gin.H{
					"Type": "System",
				},
				"ProductInfo": gin.H{
					"ID":              product.ProductID,
					"DateCreation":    (product.DateCreation()).Format("2006-01-02"),
					"DateLastBill":    (product.DateLastBill()).Format("2006-01-02"),
					"DateTermination": (product.DateTermination()).Format("2006-01-02"),
					"Terminated":      product.Terminated(),
					"WalletID":        product.WalletID,
					"BillingOption":   product.BillingOption,
				},
			})
		}

		c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, exported))
	}
}
func (Product) ListAll(c *gin.Context, user *auth.User) {
	// To list all, the user must be GLOBAL_ADMIN
	if user == nil || !user.Role.Includes(auth.GLOBAL_ADMIN) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// List products
	products, err := billing.ListAllProducts()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// For every product, build exported data
	var exported []gin.H
	for _, product := range products {
		var owner gin.H
		if product.OwnerUserID != 0 {
			owner = gin.H{
				"Type": "User",
				"ID":   product.OwnerUserID,
			}
		} else if product.OwnerAffiliationID != 0 {
			owner = gin.H{
				"Type": "Affiliation",
				"ID":   product.OwnerAffiliationID,
			}
		} else {
			owner = gin.H{
				"Type": "System",
			}
		}
		exported = append(exported, gin.H{
			"SerialNumber": product.SerialNumber(),
			"Owner":        owner,
			"ProductInfo": gin.H{
				"ID":              product.ProductID,
				"DateCreation":    (product.DateCreation()).Format("2006-01-02"),
				"DateLastBill":    (product.DateLastBill()).Format("2006-01-02"),
				"DateTermination": (product.DateTermination()).Format("2006-01-02"),
				"Terminated":      product.Terminated(),
				"WalletID":        product.WalletID,
				"BillingOption":   product.BillingOption,
			},
		})
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, exported))
}

func (Product) CreateByListingID(c *gin.Context, currentUser *auth.User) {
	// A Product may be created for
	// - a user, as PRIVATE_PRODUCT by anyone
	// - an affiliation, as AFFILIATION_PRODUCT by AFFILIATION_PRODUCT_ADMIN
	// - the system, as SYSTEM_PRODUCT by GLOBAL_ADMIN
	var form FormCreateProduct = FormCreateProduct{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if currentUser == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	// Find the Product Listing
	productListing, err := billing.GetProductListingByID(form.ProductID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var ownerUserID, ownerAffiliationID uint64
	var BillingCycle uint8 = form.BillingCycle
	var WalletID uint64

	// Finalize owner
	if currentUser.Role.Includes(auth.GLOBAL_ADMIN) {
		// Skip searching for owner ID, directly fill in from the form
		ownerUserID = form.OwnerUserOverride
		ownerAffiliationID = form.OwnerAffiliationOverride
	} else {
		// Find the owner ID
		if form.PayerOption == 0 {
			// Paid by current user
			ownerUserID = currentUser.ID()
		} else if form.PayerOption == 1 {
			// Paid by affiliation, AFFILIATION_BILLING_USER only
			if !currentUser.Role.Includes(auth.AFFILIATION_BILLING_USER) {
				c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
				return
			}
			ownerAffiliationID = currentUser.AffiliationID
		} else if form.PayerOption == 2 {
			// System product, GLOBAL_INTERNAL_USER only
			if !currentUser.Role.Includes(auth.GLOBAL_INTERNAL_USER) && !currentUser.Role.Includes(auth.GLOBAL_ADMIN) {
				c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
				return
			}
			ownerUserID = currentUser.ID()
		} else {
			// Bad
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
			return
		}
	}

	// Finalize WalletID
	if form.PayerOption == 0 {
		// Paid by the current user
		var uid uint64
		if currentUser.Role.Includes(auth.GLOBAL_ADMIN) && form.PayerOverride != 0 {
			uid = form.PayerOverride
		} else {
			uid = currentUser.ID()
		}
		wallet, err := billing.UserWallet(uid)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		WalletID = wallet.ID()
	} else if form.PayerOption == 1 {
		// Paid by affiliation
		var aid uint64
		if currentUser.Role.Includes(auth.GLOBAL_ADMIN) && form.PayerOverride != 0 {
			aid = form.PayerOverride
		} else if currentUser.Role.Includes(auth.AFFILIATION_BILLING_USER) && currentUser.AffiliationID != 0 {
			aid = currentUser.AffiliationID
		} else {
			// BAD
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
			return
		}
		affiliation, err := auth.GetAffiliationByID(aid)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		wallet, err := billing.GetWalletByID(affiliation.SharedWalletID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		WalletID = wallet.ID()
	} else if form.PayerOption == 2 { // Free
		// Must be GLOBAL_INTERNAL_USER
		if !currentUser.Role.Includes(auth.GLOBAL_INTERNAL_USER) && !currentUser.Role.Includes(auth.GLOBAL_ADMIN) {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		}
		WalletID = 0
	}

	// Create the Product
	product, err := productListing.CreateProduct(ownerUserID, ownerAffiliationID, BillingCycle, WalletID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Finalize price
	// if form.PromoCode != "" {
	// 	// Reserved for future
	// }
	if form.PriceOverride != 0 && currentUser.Role.Includes(auth.GLOBAL_ADMIN) {
		product.BillingOption.Price = form.PriceOverride
	}

	// Before creating the product in the database, charge the user if it is a recurring product
	if product.BillingOption.BillingCycle != billing.USAGE_BASED {
		wallet, err := billing.GetWalletByID(WalletID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		if err := wallet.TrySpend(product.BillingOption.Price); err != nil {
			if err == billing.ErrInsufficientFunds {
				c.AbortWithStatusJSON(http.StatusPaymentRequired, utils.RespInsufficientFunds)
			} else {
				utils.HandleError(c, err)
			}
			return
		}
	}

	// Add the product to the database
	productSN, err := product.Add()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Create Account on Server
	serverInstance, err := server.NewProvisioningServer(productListing.ServerType, productListing.ServerInstanceID, productListing.ServerConfiguration)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	err = serverInstance.CreateAccount(productSN, product.BillingOption.AccountConfiguration) // TODO: Save preset on server, not on client
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, utils.RespOK)
}
func (Product) Update(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}
func (Product) Terminate(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}
func (Product) ScheduleForTerminate(c *gin.Context, currentUser *auth.User) {
	var form FormScheduleForTerminate = FormScheduleForTerminate{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	// Parse the date
	var date time.Time
	var err error
	if form.TerminationDate != "" {
		date, err = time.Parse(time.RFC3339, form.TerminationDate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
			return
		}
	} else {
		date = time.Now()
	}

	// Get the product
	product, err := billing.GetProductBySerialNumber(form.ProductSN)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Current user must be one of the following:
	// - The user specified by product.OwnerUserID (when OwnerUserID != 0)
	// - An AFFILIATION_PRODUCT_ADMIN of the affiliation specified by GetAffiliationByID(product.OwnerAffiliationID) (when OwnerUserID == 0, OwnerAffiliationID != 0)
	// - GLOBAL_ADMIN
	if !productIsManagedByUser(product, currentUser) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	} else {
		// Set the termination date
		if err = product.ToTerminateOn(date); err != nil {
			utils.HandleError(c, err)
			return
		} else {
			c.JSON(http.StatusOK, utils.RespOK)
		}
	}
}

type Wallet struct {
}

// ?[id=<wallet_id>]
func (Wallet) View(c *gin.Context, currentUser *auth.User) {
	var walletID uint64
	var wallet *billing.Wallet
	var err error

	walletID, err = strconv.ParseUint(c.DefaultQuery("id", "0"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if walletID == 0 {
		// Get the wallet of the current user
		wallet, err = billing.UserWallet(currentUser.ID())
	} else {
		// Either of these two conditions must be met:
		// - walletID == GetAffiliationByID(currentUser.AffiliationID).SharedWalletID and currentUser.Role.Includes(auth.AFFILIATION_BILLING_ADMIN)
		// - currentUser.Role.Includes(auth.GLOBAL_ADMIN)
		if currentUser.Role.Includes(auth.GLOBAL_ADMIN) {
			wallet, err = billing.GetWalletByID(walletID)
		} else if currentUser.Role.Includes(auth.AFFILIATION_BILLING_ADMIN) {
			affiliation, err := auth.GetAffiliationByID(currentUser.AffiliationID)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			if walletID != affiliation.SharedWalletID {
				c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
				return
			}

			wallet, err = billing.GetWalletByID(walletID)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		}
	}

	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, gin.H{
		"id":       wallet.ID(),
		"owner":    wallet.OwnerUserID(),
		"balance":  wallet.Balance(),
		"secured":  wallet.Secured(),
		"disabled": wallet.Disabled(),
	}))
}

func (Wallet) Deposit(c *gin.Context, currentUser *auth.User) {
	var form FormDepositToWallet = FormDepositToWallet{}
	var wallet *billing.Wallet
	var err error

	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if form.WalletID != 0 {
		if currentUser.Role.Includes(auth.GLOBAL_ADMIN) {
			wallet, err = billing.GetWalletByID(form.WalletID)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
		} else if currentUser.Role.Includes(auth.AFFILIATION_BILLING_ADMIN) {
			affiliation, err := auth.GetAffiliationByID(currentUser.AffiliationID)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			allowedID := affiliation.SharedWalletID
			if form.WalletID != allowedID {
				c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
				return
			}
			wallet, err = billing.GetWalletByID(form.WalletID)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		}
	} else {
		wallet, err = billing.UserWallet(currentUser.ID())
		if err != nil {
			utils.HandleError(c, err)
			return
		}
	}

	// Generate checkout form
	gateway, err := payment.GetPrepaidGateway(form.PaymentInstanceID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	referenceID := walletID2ReferenceID(wallet.ID())
	paymentRequest := payment.PaymentRequest{
		Item: payment.PaymentUnit{
			ReferenceID: referenceID,
			Currency:    "USD",
			Price:       form.Amount,
		},
	}
	checkoutRenderParams, err := gateway.CheckoutForm(paymentRequest)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, gin.H{
		"deposit_to": wallet.ID(),
		"payment": gin.H{
			"gateway_instance": form.PaymentInstanceID,
			"render":           checkoutRenderParams,
		},
	}))
}
func (Wallet) Withdraw(c *gin.Context, currentUser *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}
func (Wallet) Enable(c *gin.Context, currentUser *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}
func (Wallet) Disable(c *gin.Context, currentUser *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}
