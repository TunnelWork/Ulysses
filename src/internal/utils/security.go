package utils

import (
	"crypto/ed25519"

	themis "github.com/TunnelWork/Themis"

	cbcrypt "golang.org/x/crypto/bcrypt"
)

var (
	// Authorization Token
	TokenRevoker themis.Revoker
	TokenPrivKey ed25519.PrivateKey
	TokenPubKey  ed25519.PublicKey
)

// TODO: Move HashPassword() to security package
func HashPassword(password string) string {
	hash, err := cbcrypt.GenerateFromPassword([]byte(password), cbcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
