package payment

import (
	"database/sql"
	"errors"
)

type gatewayGenFunc func(*sql.DB) (Gateway, error)

var (
	gatewayRegistry map[string]gatewayGenFunc = map[string]gatewayGenFunc{}
)

func RegisterGatewayGenFun(gatewayName string, genFunc gatewayGenFunc) {
	gatewayRegistry[gatewayName] = genFunc
}

func NewGateway(gatewayName string, db *sql.DB) (Gateway, error) {
	if genFunc, ok := gatewayRegistry[gatewayName]; ok {
		return genFunc(db)
	} else {
		return nil, errors.New("payment: gateway name not found")
	}
}
