package server

import "errors"

var (
	ErrServerConfigurables  = errors.New("BAD_SERV_CONF")
	ErrAccountConfigurables = errors.New("BAD_ACCT_CONF")
)
