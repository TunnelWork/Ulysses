package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/TunnelWork/Ulysses/src/internal/conf"
	"github.com/TunnelWork/Ulysses/src/internal/db"
	"github.com/TunnelWork/Ulysses/src/internal/logger"
	"github.com/gin-gonic/gin"
)

var (
	configPath   string
	masterConfig conf.Config
)

var (
	// Debug Switches
	noDatabase bool
	noLogger   bool
	noTick     bool
	noApi      bool

	// Global Shared Objects
	dbConnector *db.MysqlConnector
)

func parseArgs() {
	flag.StringVar(&configPath, "config", "./conf/ulysses.yaml", "Ulysses Configuration File")
	flag.BoolVar(&noDatabase, "no-db", false, "Not to use database. Thus all DB operations will give error.")
	flag.BoolVar(&noLogger, "no-log", false, "Not to use logger. Thus logging functions will do literally nothing.")
	flag.BoolVar(&noTick, "no-tick", false, "Not to use ticker. No system ticking.")
	flag.BoolVar(&noApi, "no-api", false, "Not to register API endpoints. No gin-gonic/gin ability.")
	flag.Parse()
}

func init() {
	/*** GLOBAL INIT BEGIN ***/
	parseArgs()
	globalInit()
	/*** GLOBAL INIT END ***/

	// initUlyssesServer()

}

func globalInit() {
	/*** Sync ***/
	initWaitGroup()

	/*** Internal module ***/
	initConfig()

	if !noLogger {
		initLogger()
	} else {
		fmt.Println("initLogger(): --no-log detected, skipping. What Age is this, Dark Age?")
	}
	if !noDatabase {
		initDB()
	} else {
		logger.Warning("initDB(): --no-db detected, skipping.")
	}

	/*** System module ***/
	if !noTick {
		initSystemTicking()
	} else {
		logger.Warning("initSystemTicking(): --no-tick detected, skipping.")
	}

	if !noApi {
		initApiHandler()
	} else {
		logger.Warning("initApiHandler(): --no-api detected, skipping.")
	}
}

func initWaitGroup() {
	globalWaitGroup = sync.WaitGroup{}
	globalExitChannel = make(chan bool)
}

func initConfig() {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("loadConfig(): can't read config file located at %s. error: %s", configPath, err))
	}
	masterConfig, err = conf.LoadUlyssesConfig(content)

	if err != nil {
		panic(fmt.Sprintf("loadConfig(): config file at %s is opened but can't be recognized. error: %s", configPath, err))
	}
}

// initLogger() can ONLY be called after loadConfig()
func initLogger() {
	// if err := logger.Init(masterConfig.Log); err != nil {
	var exitFunc func() = globalExitSignal
	if err := logger.InitWithWaitGroupAndExitingFunc(&globalWaitGroup, &exitFunc, masterConfig.Log); err != nil {
		panic(err)
	} else {
		logger.Info("initLogger(): success")
	}
}

// initDB() SHOULD be called after initLogger()
func initDB() {
	// DB Conn Livess Test, block until fail. No timeout.
	dbConnector = db.NewMysqlConnector(masterConfig.DB)
	if dbConnector == nil {
		logger.Fatal("initDB(): cannot establish database connection")
		return
	} else {
		logger.Info("initDB(): success")
	}
}

func initSystemTicking() {
	tickEventMutex = &sync.Mutex{}
	if masterConfig.Sys.SystemTickPeriodMillisecond == 0 {
		tickEventPeriodMs = tickEventPeriodMsDefault
	} else {
		tickEventPeriodMs = masterConfig.Sys.SystemTickPeriodMillisecond
	}
	tickEventMap = map[tickEventSignature]tickEvent{}
}

// func initUlyssesServer() {
// 	if serverConfMapMutex == nil {
// 		serverConfMapMutex = &sync.Mutex{}
// 	}
// 	serverConfMapDirty = true
// 	reloadUlyssesServer()

// 	registerTickEvent(reloadUlyssesServerSignature, reloadUlyssesServer)
// }

func initApiHandler() {
	ginRouter = gin.New()
	ginRouter.Use(gin.LoggerWithWriter(logger.NewCustomWriter("", "")), gin.Recovery())

	registerSystemAPIs()
}
