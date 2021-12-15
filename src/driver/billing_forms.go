package driver

import "github.com/TunnelWork/Ulysses.Lib/billing"

// ProductListingGroup
type (
	FormCreateProductListingGroup struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Hidden      bool   `json:"hidden"`
	}

	FormUpdateProductListingGroup struct {
		ID          uint64 `json:"id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Hidden      bool   `json:"hidden"`
	}

	FormDeleteProductListingGroup struct {
		ID uint64 `json:"id" binding:"required"`
	}
)

// ProductListing
type (
	FormProductListingID struct {
		ID uint64 `json:"id" binding:"required"`
	}

	FormCreateProductListing struct {
		ProductListing billing.ProductListing `json:"product_listing" binding:"required"`
	}

	FormUpdateProductListing struct {
		ID             uint64                 `json:"id" binding:"required"`
		ProductListing billing.ProductListing `json:"product_listing" binding:"required"`
	}

	FormDeleteProductListing struct {
		FormProductListingID
	}

	FormToggleProductListing struct {
		FormProductListingID
	}
)

// Product
type (
	FormCreateProduct struct {
		ProductID    uint64 `json:"product_id" binding:"required"`    // ID of Product Listing
		BillingCycle uint8  `json:"billing_cycle" binding:"required"` // Defined in billing/billing.go
		PayerOption  uint8  `json:"payer_option"`                     // 0 - product_owner, 1 - affiliation (AFFILIATION_BILLING_USER only), 2 - free (GLOBAL_INTERNAL_USER only)
		PromoCode    string `json:"promo_code"`                       // Promo code, reserved for future

		// GLOBAL_ADMIN overriding
		OwnerUserOverride        uint64  `json:"owner_user_override"`        // User ID of the owner of the product, default to be the Current User when PayerOption = 0. GLOBAL_ADMIN only.
		OwnerAffiliationOverride uint64  `json:"owner_affiliation_override"` // Affiliation ID of the owner of the product, default to be the Current User's Affiliation when PayerOption = 1. GLOBAL_ADMIN only.
		PayerOverride            uint64  `json:"payer_override"`             // User/Affiliation ID of the owner of the wallet to pay for the product. GLOBAL_ADMIN only.
		PriceOverride            float64 `json:"price_override"`             // Price of the product, default to be the product's price. GLOBAL_ADMIN only.

		// Product AccountConfig is saved on server with BillingOptions
		// AccountConfig interface{} `json:"account_config"` // AccountConfig of the product, used for ProvisioningServer
	}

	FormUpdateProduct struct {
		ProductSN string `json:"product_sn" binding:"required"` // Serial Number of the product
		// ProductID    uint64 `json:"product_id"` // Should not be updated
		// BillingCycle uint8  `json:"billing_cycle"` // Should not be updated
		PayerOption uint8  `json:"payer_option"` // 0 - product_owner, 1 - affiliation (AFFILIATION_BILLING_USER only), 2 - free (GLOBAL_INTERNAL_USER only)
		PromoCode   string `json:"promo_code"`   // Promo code, reserved for future

		// GLOBAL_ADMIN overriding
		OwnerUser        uint64  `json:"owner_user_override"`        // User ID of the owner of the product, default to be the Current User when PayerOption = 0. GLOBAL_ADMIN only.
		OwnerAffiliation uint64  `json:"owner_affiliation_override"` // Affiliation ID of the owner of the product, default to be the Current User's Affiliation when PayerOption = 1. GLOBAL_ADMIN only.
		PayerOverride    uint64  `json:"payer_override"`             // User/Affiliation ID of the owner of the wallet to pay for the product. GLOBAL_ADMIN only.
		PriceOverride    float64 `json:"price_override"`             // Price of the product, default to be the product's price. GLOBAL_ADMIN only.

		// Product AccountConfig
		// AccountConfig interface{} `json:"account_config"` // Should not be updated
	}

	FormScheduleForTerminate struct {
		ProductSN       uint64 `json:"product_sn" binding:"required"` // Serial Number of the product
		TerminationDate string `json:"termination_date"`              // Format: YYYY-MM-DD
	}
)

// Wallet
type (
	FormDepositToWallet struct {
		WalletID          uint64  `json:"wallet_id"`                              // default: user's wallet
		Amount            float64 `json:"amount" binding:"required"`              // Amount, in USD, to deposit
		PaymentInstanceID string  `json:"payment_instance_id" binding:"required"` // Payment method
	}
)
