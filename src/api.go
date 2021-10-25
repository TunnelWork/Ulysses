package main

import (
	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/TunnelWork/Ulysses/src/internal/logger"
	"github.com/gin-gonic/gin"
)

var (
	// MUST REGISTER ALL FUNCTION HERE
	mapSystemApiPostHandlers = map[string](*gin.HandlerFunc){
		"MFA":  &handlerCheckMFA,
		"Auth": &handlerAuth,
	}

	mapSystemApiGetHandlers = map[string](*gin.HandlerFunc){}

	mapDebugPost = map[string](*gin.HandlerFunc){
		"debug/SM": &handlerDebugSM,
	}
	mapDebugGet = map[string](*gin.HandlerFunc){}

	// MUST CREATE FUNCTION VARIABLE AS POINTER
	handlerCheckMFA gin.HandlerFunc = _handlerCheckMFA
	handlerAuth     gin.HandlerFunc = _handlerAuth
	handlerDebugSM  gin.HandlerFunc = _debugHandlerServerManager
)

// registerSystemAPIs() is just an additional step to prevent API endpoints confliction.
// it reuse the register route a third-party module will use.
func registerSystemAPIs() {
	var err error
	for route, handler := range mapSystemApiPostHandlers {
		err = api.RegisterApiEndpoint(api.HTTP_METHOD_POST, route, handler)
		if err != nil {
			logger.Fatal("registerSystemAPIs(): Cannot register POST route", route, " due to error: ", err)
		}
	}

	for route, handler := range mapSystemApiGetHandlers {
		api.RegisterApiEndpoint(api.HTTP_METHOD_GET, route, handler)
		if err != nil {
			logger.Fatal("registerSystemAPIs(): Cannot register GET route", route, " due to error: ", err)
		}
	}

	for route, handler := range mapDebugGet {
		err = api.RegisterApiEndpoint(api.HTTP_METHOD_GET, route, handler)
		if err != nil {
			logger.Fatal("registerSystemAPIs(): Cannot register POST route", route, " due to error: ", err)
		}
	}

	for route, handler := range mapDebugPost {
		err = api.RegisterApiEndpoint(api.HTTP_METHOD_POST, route, handler)
		if err != nil {
			logger.Fatal("registerSystemAPIs(): Cannot register POST route", route, " due to error: ", err)
		}
	}

	api.ImportToGinEngine(ginRouter, masterConfig.Sys.UrlPath)
}
