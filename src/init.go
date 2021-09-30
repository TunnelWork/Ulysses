package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/TunnelWork/Ulysses/src/internal/conf"
	"github.com/TunnelWork/Ulysses/src/internal/db"
	"github.com/TunnelWork/Ulysses/src/internal/logger"
)

// Logging Level
// 0: Absolute No Logging, you don't even know how it crashed (unless panic)
// 5: Log every non-trivial thing. FBI Open Up!
const (
	logNull uint8 = iota
	logFatal
	logError
	logWarning
	logInfo
	logDebug
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

// initLogger() can ONLY be called after loadConfig()
func initLogger() {
	if err := logger.Init(masterConfig.Log); err != nil {
		panic(err)
	}
}

// initDB() SHOULD be called after initLogger()
func initDB() {
	// DB Conn Livess Test
	_, err := db.DBConnect(masterConfig.DB)
	if err != nil {
		logger.Fatal("initDB(): ", err)
	} else {
		logger.Debug("initDB(): PASS")
	}
}

func init() {
	//// GLOBAL INIT BEGIN ////
	flag.StringVar(&configPath, "configPath", "./conf/ulysses.yaml", "Ulysses Configuration File")
	flag.Parse()

	loadConfig()
	/*** GLOBAL INIT END ***/
	/*** SYSTEM MODULE INIT BEGIN ***/
	initLogger() // First thing first or we have nowhere to write
	initDB()
	/*** SYSTEM MODULE INIT END ***/

	select {}
}
