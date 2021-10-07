package api

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mapMutex sync.RWMutex = sync.RWMutex{}

	apiGETMap  map[string]*gin.HandlerFunc
	apiPOSTMap map[string]*gin.HandlerFunc
)

func ImportToGinEngine(router *gin.Engine) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	for path, handler := range apiGETMap {
		router.GET(path, *handler)
	}

	for path, handler := range apiPOSTMap {
		router.POST(path, *handler)
	}
}

func ExportHandlerMaps() (getMap map[string]*gin.HandlerFunc, postMap map[string]*gin.HandlerFunc) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	// Copy list
	getMap = apiGETMap
	postMap = apiPOSTMap

	return getMap, postMap
}
