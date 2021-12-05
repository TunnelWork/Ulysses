package main

import (
	"database/sql"
	"os"
	"sync"

	uconf "github.com/TunnelWork/Ulysses/src/internal/conf"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

var (
	// Config
	ulyssesConfigPath string
	ulyssesConfigFile uconf.UlyssesConfigFile
	completeConfig    uconf.CompleteConfig

	// Database
	dbPool *sql.DB

	// Cron
	crontab *cron.Cron

	// System/Sync
	sysSig chan os.Signal = make(chan os.Signal, 1)

	// // Auth: MFA Plugin
	// totp     *utotp.TOTP
	// webauthn *uwebauthn.WebAuthn

	// // Auth: Token verification
	// revoker      *themis.OfflineRevoker
	// tokenPrivKey ed25519.PrivateKey
	// tokenPubKey  ed25519.PublicKey

	// Gin/API
	ginRouter *gin.Engine = gin.Default()

	// Sync Utils
	masterWaitGroup   sync.WaitGroup // main() will Wait() on this group before exiting. Slaves must Add(1) on this before Business Logic and Done() after Business Logic.
	slaveWaitGroup    sync.WaitGroup // main() will Add(1) to this in exiting procedual. Slaves must Wait(), i.e., not to start Business Logic once exiting procedual has started.
	globalExitChannel chan bool
)
