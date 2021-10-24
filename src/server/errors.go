package server

import "errors"

// Caller of this module should check the returned error against the list provided below.
// Any error returned not appearing in the list should be considered as an Internal Error
// and should not be displayed to Customer/Non-Dev.
var (
	ErrServerUnknown        = errors.New("ulysses/server: server type is not registered")
	ErrServerConfigurables  = errors.New("ulysses/server: bad server config")
	ErrAccountConfigurables = errors.New("ulysses/server: bad account config")
	ErrBadJsonObject        = errors.New("ulysses/server: bad JSON object")
	ErrBadJsonArray         = errors.New("ulysses/server: bad JSON array")
)
