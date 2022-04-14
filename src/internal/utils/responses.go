package utils

import "github.com/TunnelWork/Ulysses.Lib/api"

var (
	// General
	RespOK = api.MessageResponse(api.SUCCESS, "SUCCESS")

	RespBadRequest     = api.MessageResponse(api.ERROR, "BAD_REQUEST")
	RespNotAuthorized  = api.MessageResponse(api.ERROR, "NOT_AUTHORIZED") // invalid credential
	RespAccessDenied   = api.MessageResponse(api.ERROR, "ACCESS_DENIED")  // valid credential, not authorized for resource
	RespInternalErr    = api.MessageResponse(api.ERROR, "INTERNAL_SERVER_ERROR")
	RespNotImplemented = api.MessageResponse(api.ERROR, "NOT_IMPLEMENTED")

	// Authorization
	RespAccessTokenRequired = api.MessageResponse(api.ERROR, "ACCESS_TOKEN_REQUIRED")
	RespAccessTokenInvalid  = api.MessageResponse(api.ERROR, "ACCESS_TOKEN_INVALID")

	// User
	RespUserNotFound   = api.MessageResponse(api.ERROR, "USER_NOT_FOUND")
	RespUserInfoNeeded = api.MessageResponse(api.SUCCESS, "USER_INFO_NEEDED")

	// MFA
	RespMfaNotFound         = api.MessageResponse(api.ERROR, "MFA_NOT_FOUND")
	RespMfaRequestInvalid   = api.MessageResponse(api.ERROR, "MFA_REQUEST_INVALID")
	RespMfaResponseRequired = api.MessageResponse(api.ERROR, "MFA_RESPONSE_REQUIRED")
	RespMfaResponseInvalid  = api.MessageResponse(api.ERROR, "MFA_RESPONSE_INVALID")

	// Database
	RespBadDatabase = api.MessageResponse(api.ERROR, "ERR_BAD_DATABASE")

	// Email
	RespInvalidEmail  = api.MessageResponse(api.ERROR, "ERR_EMAIL_INVALID")
	RespEmailConflict = api.MessageResponse(api.ERROR, "ERR_EMAIL_CONFLICT")

	// Billing
	RespInsufficientFunds = api.MessageResponse(api.ERROR, "ERR_INSUFFICIENT_FUNDS")
)
