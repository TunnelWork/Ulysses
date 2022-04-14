package main

import (
	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/gin-gonic/gin"
)

type handler = gin.HandlerFunc

var (
	GETAuth map[string][]*handler = map[string][]*handler{
		"affiliation": {
			&AuthorizationMustBeValid,
			&GETAuthAffiliation,
		},
		"mfa": {
			&AuthorizationMustBeValid,
			&GETAuthMFA,
		},
		"user": {
			&AuthorizationMustBeValid,
			&GETAuthUser,
		},
	}
	POSTAuth map[string][]*handler = map[string][]*handler{
		"affiliation": {
			&AuthorizationMustBeValid,
			&POSTAuthAffiliation,
		},
		"mfa": {
			&AuthorizationMustBeValid,
			&POSTAuthMFA,
		},
		"user": {
			&AuthorizationMustBeValidIfExists,
			&POSTAuthUser,
		},
	}

	GETBilling map[string][]*handler = map[string][]*handler{
		"billingrecord": {
			&AuthorizationMustBeValid,
			&GETBillingRecord,
		},
		"productlistinggroup": {
			&GETBillingProductListingGroup,
		},
		"productlisting": {
			&AuthorizationMustBeValidIfExists,
			&POSTBillingProductListing,
		},
		"product": {
			&AuthorizationMustBeValid,
			&GETBillingProduct,
		},
		"wallet": {
			&AuthorizationMustBeValid,
			&GETBillingWallet,
		},
	}
	POSTBilling map[string][]*handler = map[string][]*handler{
		"billingrecord": {
			&AuthorizationMustBeValid,
			&POSTBillingRecord,
		},
		"productlistinggroup": {
			&AuthorizationMustBeValid,
			&UserMustBeGlobalAdmin,
			&MFAMustBeEnabled,
			&MFARespMustBeValid,
			&POSTBillingProductListingGroup,
		},
		"productlisting": {
			&AuthorizationMustBeValid,
			&UserMustBeGlobalAdmin,
			&MFAMustBeEnabled,
			&MFARespMustBeValid,
			&POSTBillingProductListing,
		},
		"product": {
			&AuthorizationMustBeValid,
			&MFAMustBeEnabled,
			&MFARespMustBeValid,
			&POSTBillingProduct,
		},
		"wallet": {
			&AuthorizationMustBeValid,
			&MFAMustBeEnabled,
			&MFARespMustBeValid,
			&POSTBillingWallet,
		},
	}

	GETServer map[string][]*handler = map[string][]*handler{
		"provisioning/account": {
			&AuthorizationMustBeValid,
			&GETProvisioningAccount,
		},
	}
	POSTServer map[string][]*handler = map[string][]*handler{
		"provisioning/account": {
			&AuthorizationMustBeValid,
			&UserMustBeGlobalAdmin,
			&MFAMustBeEnabled,
			&MFARespMustBeValid,
			&POSTProvisioningAccount,
		},
	}
)

// When checkpoint/endpoint fails, it always respond with api.MessageResponse
// When checkpoint/endpoint success, it may respond with api.MessageResponse or api.PayloadResponse
func registerAPIEndpoints() error {
	var err error

	// Authorize
	err = api.POST("authorize", &Authorize)
	if err != nil {
		return err
	}

	// Auth
	for path, handlers := range GETAuth {
		err = api.CGET(api.Auth, path, handlers...)
		if err != nil {
			return err
		}
	}
	for path, handlers := range POSTAuth {
		err = api.CPOST(api.Auth, path, handlers...)
		if err != nil {
			return err
		}
	}

	// Billing
	for path, handlers := range GETBilling {
		err = api.CGET(api.Billing, path, handlers...)
		if err != nil {
			return err
		}
	}
	for path, handlers := range POSTBilling {
		err = api.CPOST(api.Billing, path, handlers...)
		if err != nil {
			return err
		}
	}

	// Server
	for path, handlers := range GETServer {
		err = api.CGET(api.Server, path, handlers...)
		if err != nil {
			return err
		}
	}
	for path, handlers := range POSTServer {
		err = api.CPOST(api.Server, path, handlers...)
		if err != nil {
			return err
		}
	}
	return nil
}
