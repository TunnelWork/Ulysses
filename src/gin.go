package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ginRouter *gin.Engine
)

func startGinRouter() {
	ginRouter.Run(fmt.Sprintf("%s:%d", masterConfig.Sys.Host, masterConfig.Sys.Port))
}
