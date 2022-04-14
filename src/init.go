package main

import (
	"crypto/ed25519"
	"flag"
	"fmt"
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

	paypal "github.com/TunnelWork/payment.PayPal/v2"
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
	// Everyday at midnight
	crontab.AddFunc("0 0 0 * * *", func() {
		errs := billing.DailyRecurringBilling()
		if len(errs) > 0 {
			logging.Error("Cronjob: DailyRecurringBilling() returned error(s):")
			for _, err := range errs {
				logging.Error(err.Error())
			}
		} else {
			logging.Info("Cronjob: DailyRecurringBilling() executed successfully.")
		}
	})
	// Every hour, sync usage-based billing
	crontab.AddFunc("0 0 * * * *", func() {
		errs := billing.HourlyUsageBilling()
		if len(errs) > 0 {
			logging.Error("Cronjob: HourlyUsageBilling() returned error(s):")
			for _, err := range errs {
				logging.Error(err.Error())
			}
		} else {
			logging.Info("Cronjob: HourlyUsageBilling() executed successfully.")
		}
	})
	// Every hour, terminate products that is
	crontab.AddFunc("0 0 * * * *", func() {
		errs := billing.HourlyProductTermination()
		if len(errs) > 0 {
			logging.Error("Cronjob: HourlyProductTermination() returned error(s):")
			for _, err := range errs {
				logging.Error(err.Error())
			}
		} else {
			logging.Info("Cronjob: HourlyProductTermination() executed successfully.")
		}
	})
	// Every hour, purge expired tmp database entries
	crontab.AddFunc("0 0 * * * *", func() {
		auth.PurgeExpiredTmpEntry()
		logging.Debug("Cronjob: Purged Tmp Database.")
	})

	/*** Rest of the Ulysses.Lib ***/
	auth.Setup(dbPool, completeConfig.Mysql.TblPrefix)
	billing.Setup(dbPool, completeConfig.Mysql.TblPrefix)
	payment.Setup(dbPool, completeConfig.Mysql.TblPrefix) // TODO
	security.SetupCipher(completeConfig.Security.Cipher())

	/*** Plugins ***/
	// MFA
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
	// Payment
	if paypalPrepaidConfig, err := paypal.LoadPrepaidConfig(dbPool, completeConfig.Mysql.TblPrefix, "paypal_wallet"); err != nil {
		logging.Fatal("cannot load config for paypal prepaid gateway, error: ", err)
	} else {
		_, err := payment.NewPrepaidGateway("paypal", "paypal_prepaid_wallet_deposit", map[string]string{
			"clientID":     paypalPrepaidConfig.ClientID,
			"secretID":     paypalPrepaidConfig.SecretID,
			"apiBase":      paypalPrepaidConfig.ApiBase,
			"callbackBase": fmt.Sprintf("%s/%spayment/callback/", completeConfig.Http.URLComplete, completeConfig.Http.URLPrefix),
		})
		if err != nil {
			logging.Fatal("cannot initialize paypal prepaid gateway, error: ", err)
		}
	}

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
