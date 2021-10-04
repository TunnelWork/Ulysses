package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/TunnelWork/Ulysses/src/internal/conf"
	"github.com/TunnelWork/Ulysses/src/internal/db"
	"github.com/TunnelWork/Ulysses/src/internal/logger"
)

var (
	configPath   string
	masterConfig conf.Config
)

func loadConfig() {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("loadConfig(): can't read config file located at %s. error: %s", configPath, err))
	}
	masterConfig, err = conf.LoadUlyssesConfig(content)

	if err != nil {
		panic(fmt.Sprintf("loadConfig(): config file at %s is opened but can't be recognized. error: %s", configPath, err))
	}
}

var (
	// Debug Switches
	noDatabase bool
	noLogger   bool
	noTick     bool
)

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
	_, err := db.DBConnect(masterConfig.DB)
	if err != nil {
		logger.Fatal("initDB(): ", err)
		return
	} else {
		logger.Info("initDB(): success")
	}
}

func parseArgs() {
	flag.StringVar(&configPath, "config", "./conf/ulysses.yaml", "Ulysses Configuration File")
	flag.BoolVar(&noDatabase, "no-db", false, "Not to use database. Thus all DB operations will give error.")
	flag.BoolVar(&noLogger, "no-log", false, "Not to use logger. Thus logging functions will do literally nothing.")
	flag.BoolVar(&noTick, "no-tick", false, "Not to use ticker. No system ticking.")
	flag.Parse()
}

func globalInit() {
	loadConfig()
	initWaitGroup()
}

func init() {
	/*** GLOBAL INIT BEGIN ***/
	parseArgs()
	globalInit()
	/*** GLOBAL INIT END ***/

	/*** INTERNAL MODULE INIT BEGIN ***/
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
	/*** INTERNAL MODULE INIT END ***/

	/*** SYSTEM MODULE INIT BEGIN ***/
	if !noTick {
		initSystemTicking() // First one!
	} else {
		logger.Warning("initSystemTicking(): --no-tick detected, skipping.")
	}
	initUlyssesServer()
	/*** SYSTEM MODULE INIT END ***/
}
