package payment

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type GatewayStatusFlag uint32

const (
	GATEWAY_NEEDCONFIG GatewayStatusFlag = 1 << iota // If set, the gateway is properly configured. Otherwise it needs to be configured before use.
	GATEWAY_BADCONFIG
	GATEWAY_GOOD
)

type Gateway interface {
	/*** Status ***/
	GatewayStatus() GatewayStatusFlag

	/*** Gateway Config ***/
	// GatewayConfigTemplate() returns needed parameters in setting up the Gateway.
	// Please see examples.go for an example
	// this is used for frontend implementation compatibilities
	GatewayConfigTemplate() (gatewayConfigTemplate P)

	// GatewayConfig() saves gatewayConfiguration to database
	GatewayConfig(db *sql.DB, gatewayConfig P) error

	/*** Order Operations ***/
	// Customer(payer)-only
	CreateOrder(db *sql.DB, orderCreationParams P) (orderID string, err error)

	// Admin/Customer
	LookupOrder(db *sql.DB, ReferenceID string) (orderDetails P, err error)

	// Admin(payee)-only, should be saved when the order is created
	OrderDetail(db *sql.DB, orderID string) (orderDetails P, err error)

	// For Admin(payee)/Customer(payer), should be fetched from remote API
	// everytime gets requested
	OrderStatus(db *sql.DB, orderID string) (orderStatus P, err error)

	//
	GenerateOrderForm(db *sql.DB, orderID string) (orderFormTemplate P, err error)
	FinalizeOnSiteOrderForm(db *sql.DB, onSiteOrderForm P) error

	CancelOrder(db *sql.DB, orderID string) error
}

type GatewayCallback interface {
	Gateway

	// Set which handlerFunc() to be called with latest orderStatus for an order
	// when there's a callback from the Payment provider.
	SetCallbackHandler(handlerFunc func(orderStatus P))

	// callback URL should be a complete URL
	// e.g., https://ulysses.tunnel.work/api/callback/paypal
	// caller is responsible to bind Callback() to the same endpoint
	SetCallbackURL(callback string)

	// Actual Callback function that needs to be connected to gin Router
	Callback(db *sql.DB, c *gin.Context)
}

type GatewayRefundable interface {
	Gateway

	// Refund() the orderID in full or partial, depending on the params
	Refund(db *sql.DB, orderRefundParams P) error
}
