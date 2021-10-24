package api

import (
	"errors"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	HTTP_METHOD_GET uint = iota
	HTTP_METHOD_POST
)

var (
	ErrNotAllowDirectFuncReg error = errors.New("api.RegisterApiEndpoint(): no direct handler func registration is allowed")
	ErrBadMethod             error = errors.New("api.RegisterApiEndpoint(): bad method")
	ErrRepeatGetPath         error = errors.New("api.RegisterApiEndpoint(): repeated path for GET method")
	ErrRepeatPostPath        error = errors.New("api.RegisterApiEndpoint(): repeated path for POST method")
)

// RegisterNewAPIEndpoint() allows only one single handler function,
// unlike with direct Gin.Router access.
// note: this is overpowered and makes the system vulnerable to malicious code.
// in this version, I made it callable ONLY from main package.
func RegisterApiEndpoint(method uint, relativePath string, handler *gin.HandlerFunc) error {
	var packageName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		packageName = strings.Split(details.Name(), ".")[0]
	}

	if packageName == "main" {
		return registerApiEndpoint(method, relativePath, handler)
	} else {
		return ErrNotAllowDirectFuncReg
	}
}

func registerApiEndpoint(method uint, relativePath string, handler *gin.HandlerFunc) error {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	switch method {
	case HTTP_METHOD_GET:
		return registerApiGETEndpoint(relativePath, handler)
	case HTTP_METHOD_POST:
		return registerApiPOSTEndpoint(relativePath, handler)
	default:
		return ErrBadMethod
	}
}

func registerApiGETEndpoint(relativePath string, handler *gin.HandlerFunc) error {
	if _, ok := apiGETMap[relativePath]; ok {
		return ErrRepeatGetPath
	}
	apiGETMap[relativePath] = handler
	return nil
}

func registerApiPOSTEndpoint(relativePath string, handler *gin.HandlerFunc) error {
	if _, ok := apiPOSTMap[relativePath]; ok {
		return ErrRepeatPostPath
	}
	apiPOSTMap[relativePath] = handler
	return nil
}
