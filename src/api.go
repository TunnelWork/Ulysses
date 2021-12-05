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
)

// When checkpoint/endpoint fails, it always respond with api.MessageResponse
// When checkpoint/endpoint success, it may respond with api.MessageResponse or api.PayloadResponse
func registerAPIEndpoints() error {
	var err error
	// Auth
	for endpoint, handlers := range GETAuth {
		err = api.CGET(api.Auth, endpoint, handlers...)
		if err != nil {
			return err
		}
	}
	for endpoint, handlers := range POSTAuth {
		err = api.CPOST(api.Auth, endpoint, handlers...)
		if err != nil {
			return err
		}
	}
	return nil
}
