package server

import "errors"

var (
	ErrServerConfigurables  = errors.New("ulysses/server: bad server config")
	ErrAccountConfigurables = errors.New("ulysses/server: bad account config")
)
