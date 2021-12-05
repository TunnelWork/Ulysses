package utils

import (
	"net/http"

	harpocrates "github.com/TunnelWork/Harpocrates"
	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/TunnelWork/Ulysses.Lib/logging"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return // won't handle nil
	}
	// Generate error random index for logging
	index, err := harpocrates.GetRandomBase32()
	if err != nil {
		index = "N/A"
	}
	// Log error
	logging.Error("[%s] %s", index, err.Error())

	c.AbortWithStatusJSON(http.StatusInternalServerError, api.PayloadResponse(api.ERROR, gin.H{
		"err_index": index,
	}))
}
