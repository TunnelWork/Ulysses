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

	// Debug Switches
	skipDatabase bool
	skipLogger   bool
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
	} else {
		logger.Info("initLogger(): success")
	}
}

// initDB() SHOULD be called after initLogger()
func initDB() {
	// DB Conn Livess Test
	_, err := db.DBConnect(masterConfig.DB)
	if err != nil {
		logger.Fatal("initDB(): ", err)
	} else {
		logger.Info("initDB(): success")
	}
}

func parseArgs() {
	flag.StringVar(&configPath, "config", "./conf/ulysses.yaml", "Ulysses Configuration File")
	flag.BoolVar(&skipDatabase, "skip-db", false, "Not to use database. Thus all DB operations will give error.")
	flag.BoolVar(&skipLogger, "skip-log", false, "Not to use logger. Thus logging functions will do literally nothing.")
	flag.Parse()
}

func globalInit() {
	loadConfig()
}

func init() {
	/*** GLOBAL INIT BEGIN ***/
	parseArgs()
	globalInit()
	/*** GLOBAL INIT END ***/

	/*** SYSTEM MODULE INIT BEGIN ***/
	if !skipLogger {
		initLogger()
	} else {
		fmt.Println("initLogger(): --skip-log detected, skipping. What Age is this, Dark Age?")
	}
	if !skipDatabase {
		initDB()
	} else {
		logger.Warning("initDB(): --skip-db detected, skipping.")
	}
	/*** SYSTEM MODULE INIT END ***/
}
