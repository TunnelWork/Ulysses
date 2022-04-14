package utils

import (
	"net/http"
	"runtime"

	harpocrates "github.com/TunnelWork/Harpocrates"
	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/TunnelWork/Ulysses.Lib/logging"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return // won't handle nil
	}

	// Get information about where the error occured
	// Credits: https://stackoverflow.com/a/53012754/10469909
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	// Generate error random index for logging
	index, err := harpocrates.GetRandomBase32()
	if err != nil {
		index = "N/A"
	}
	// Log error
	logging.Error("[%s] %s: %s", index, funcName, err.Error())

	c.AbortWithStatusJSON(http.StatusInternalServerError, api.PayloadResponse(api.ERROR, gin.H{
		"ref": index,
	}))
}
