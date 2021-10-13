package api

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mapMutex sync.RWMutex = sync.RWMutex{}

	apiGETMap  map[string]*gin.HandlerFunc = map[string]*gin.HandlerFunc{}
	apiPOSTMap map[string]*gin.HandlerFunc = map[string]*gin.HandlerFunc{}
)

func ImportToGinEngine(router *gin.Engine, urlPath string) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	// Safety measure: trim-off all leading/ending slashes (/)
	for urlPath[0] == '/' {
		urlPath = urlPath[1:] // trim off the first character
	}
	for urlPath[len(urlPath)-1] == '/' {
		urlPath = urlPath[:len(urlPath)-1] // trim off the last character
	}

	// For non empty urlPath, append ending slash to make it a path.
	if len(urlPath) > 0 {
		urlPath = urlPath + "/"
	}

	for path, handler := range apiGETMap {
		router.GET(urlPath+path, *handler)
	}

	for path, handler := range apiPOSTMap {
		router.POST(urlPath+path, *handler)
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
