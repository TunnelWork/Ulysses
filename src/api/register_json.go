package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFuncJSON func(*LessContext) (status int, json gin.H)

// TODO: Add more safe-to-expose access
type LessContext struct {
	Request *http.Request
}

// This function should wrap the JSON and StatCode returned from handler, register it
func RegisterApiEndpointJSON(method uint, relativePath string, handler *HandlerFuncJSON) error {
	var wrappedHandlerFunc gin.HandlerFunc = func(c *gin.Context) {
		lc := &LessContext{
			Request: c.Request,
		}
		status, json := (*handler)(lc)
		c.JSON(status, json)
	}

	return registerApiEndpoint(method, relativePath, &wrappedHandlerFunc)
}
