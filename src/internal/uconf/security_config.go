package uconf

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"

	harpocrates "github.com/TunnelWork/Harpocrates"
	"github.com/TunnelWork/Ulysses.Lib/security"
	"golang.org/x/crypto/hkdf"
)

type SecurityModuleConfig struct {
	secSeed string
}

func (smc SecurityModuleConfig) Cipher() security.Cipher {
	if len(smc.secSeed) < 64 {
		panic("uconf: security_seed must be longer than 64 bytes")
	}

	hkdfReader := hkdf.New(sha256.New, []byte(smc.secSeed), nil, []byte("UlyssesUlyssesUlyssesUlyssesUlyssesUlysses"))

	// Read 32-byte for AES Key
	aesKey := make([]byte, 32)
	_, err := hkdfReader.Read(aesKey)
	if err != nil {
		panic(fmt.Sprintf("uconf: cannot read from hkdfReader, err: %s", err))
	}

	// Read 32-byte for AES IV
	aesIV := make([]byte, 32)
	_, err = hkdfReader.Read(aesIV)
	if err != nil {
		panic(fmt.Sprintf("uconf: cannot read from hkdfReader, err: %s", err))
	}

	// Create AES Cipher
	cipher := harpocrates.NewAESCipher(aesKey, aesIV, harpocrates.CBC)
	if cipher == nil {
		panic("uconf: cannot create AES cipher")
	}

	return cipher
}

func (smc SecurityModuleConfig) Ed25519Key() ed25519.PrivateKey {
	return harpocrates.Ed25519Key(smc.secSeed)
}
