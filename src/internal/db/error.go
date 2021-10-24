package db

import (
	"errors"
)

var (
	ErrCannotAppendCert = errors.New("/src/internal/db: AppendCertsFromPEM() failed")
	ErrIncompleteConf   = errors.New("/src/internal/db: configuration isn't complete")
	ErrMySQLNoConn      = errors.New("/src/internal/db: cannot establish mysql connection")
)
