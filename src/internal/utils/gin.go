package utils

import (
	"errors"

	themis "github.com/TunnelWork/Themis"
	"github.com/gin-gonic/gin"

	"github.com/TunnelWork/Ulysses.Lib/auth"
)

var (
	ErrNoAuthorizationHeader  = errors.New("utils: no authorization header")
	ErrBadAuthorizationHeader = errors.New("utils: bad authorization header")
)

// getUserID extracts the user ID from the Authorization header
// currently, not verifying the Authorization's validity.
func AuthorizationToUserID(c *gin.Context) (uint64, error) {
	authHeader := c.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		return 0, ErrNoAuthorizationHeader
	}
	if len(authHeader) > 1 {
		return 0, ErrBadAuthorizationHeader
	}
	authToken := authHeader[0]
	bearer, err := themis.ImportBearerToken(authToken, nil)
	if err != nil {
		return 0, err
	}
	// Commented the bearer token validation because it is TOO SLOW!
	// if err = bearer.Verify(tokenPubKey); err != nil {
	// 	return 0, err
	// }
	return bearer.Body().Identity, nil
}

// getUserID extracts the user ID from the Authorization header
// currently, not verifying the Authorization's validity.
func AuthorizationToUser(c *gin.Context) (*auth.User, error) {
	userID, err := AuthorizationToUserID(c)
	if err != nil {
		if err == ErrNoAuthorizationHeader {
			return nil, nil // Allow requesting nil user
		} else {
			return nil, err
		}
	}
	return auth.GetUserByID(userID)
}
