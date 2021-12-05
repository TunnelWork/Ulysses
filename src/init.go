package main

import (
	"crypto/ed25519"
	"flag"
	"os/signal"
	"sync"
	"syscall"

	themis "github.com/TunnelWork/Themis"
	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/TunnelWork/Ulysses.Lib/auth"
	utotp "github.com/TunnelWork/Ulysses.Lib/auth/mfa/totp"
	uwebauthn "github.com/TunnelWork/Ulysses.Lib/auth/mfa/webauthn"
	"github.com/TunnelWork/Ulysses.Lib/billing"
	"github.com/TunnelWork/Ulysses.Lib/logging"
	"github.com/TunnelWork/Ulysses.Lib/payment"
	"github.com/TunnelWork/Ulysses.Lib/security"
	"github.com/TunnelWork/Ulysses/src/internal/uconf"
	"github.com/TunnelWork/Ulysses/src/internal/utils"
	"github.com/robfig/cron/v3"

	_ "github.com/TunnelWork/payment.PayPal"
	_ "github.com/TunnelWork/server.Trojan"
)

func init() {
	parseArgs()
	globalInit()
}

func parseArgs() {
	flag.StringVar(&ulyssesConfigPath, "config", "./conf/ulysses.yaml", "Ulysses Database Configuration File")
	flag.Parse()
}

func globalInit() {
	var err error
	/*** Low-level Functionalities ***/
	masterWaitGroup = sync.WaitGroup{}
	slaveWaitGroup = sync.WaitGroup{}
	globalExitChannel = make(chan bool)
	signal.Notify(sysSig, syscall.SIGTERM, syscall.SIGINT)

	/*** Configuration ***/
	ulyssesConfigFile, err = uconf.LoadConfigFromFile(ulyssesConfigPath)
	if err != nil {
		panic(err)
	}
	completeConfig, err = ulyssesConfigFile.LoadCompleteConfig()
	if err != nil {
		panic(err)
	}
	dbPool, err = completeConfig.Mysql.DB()
	if err != nil {
		panic(err)
	}

	/*** Logger ***/
	if err := logging.InitWithWaitGroupAndExitingFunc(&masterWaitGroup, nil, completeConfig.Logger); err != nil {
		panic(err)
	} else {
		logging.Info("initLogger(): success")
	}

	/*** Cron ***/
	crontab = cron.New()
	// // Everyday at midnight
	// crontab.AddFunc("0 0 0 * * *", func() {
	// 	errs := billing.DailyRecurringBilling()
	// 	if len(errs) > 0 {
	// 		for _, err := range errs {
	// 			logging.Error(err)
	// 		}
	// 	}
	// })
	// // Every hour
	// crontab.AddFunc("0 0 * * * *", func() {
	// 	errs := billing.HourlyUsageBilling()
	// 	if len(errs) > 0 {
	// 		for _, err := range errs {
	// 			logging.Error(err)
	// 		}
	// 	}
	// })
	// // Every hour
	// crontab.AddFunc("0 0 * * * *", func() {
	// 	errs := billing.HourlyProductTermination()
	// 	if len(errs) > 0 {
	// 		for _, err := range errs {
	// 			logging.Error(err)
	// 		}
	// 	}
	// })

	/*** Rest of the Ulysses.Lib ***/
	auth.Setup(dbPool, completeConfig.Mysql.TblPrefix)
	billing.Setup(dbPool, completeConfig.Mysql.TblPrefix)
	payment.Setup(dbPool, completeConfig.Mysql.TblPrefix) // TODO
	security.SetupCipher(completeConfig.Security.Cipher())

	/*** Non-Volatile Plugins ***/
	totp := utotp.NewTOTP(map[string]string{
		"issuer": "Ulysses",
	})
	webauthn := uwebauthn.NewWebAuthn(map[string]string{
		"RPDisplayName": "Ulysses",
		"RPID":          completeConfig.Http.URLDomain,
		"RPOriginURL":   completeConfig.Http.URLComplete,
	})
	auth.RegMFAInstance("utotp", totp)
	auth.RegMFAInstance("uwebauthn", webauthn)

	utils.TokenRevoker = themis.NewOfflineRevoker()
	utils.TokenPrivKey = completeConfig.Security.Ed25519Key()
	utils.TokenPubKey = utils.TokenPrivKey.Public().(ed25519.PublicKey)

	/*** API ***/
	err = registerAPIEndpoints()
	if err != nil {
		panic(err)
	}

	api.FinalizeGinEngine(ginRouter, completeConfig.Http.URLPrefix)
	ginRouter.HandleMethodNotAllowed = true
	ginRouter.RedirectTrailingSlash = true
	ginRouter.RedirectFixedPath = true

	// TODO: Set XFF trusted proxies
}
