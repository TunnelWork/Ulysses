package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	verboseLogging bool
	logFile        string
	logLevel       uint8
}

func initConfig() {
	content, err := ioutil.ReadFile(ulyssesConfigPath)
	if err != nil {
		Fatal("initDB(): ", err)
	}

	yaml.Unmarshal(content, &ulyssesConfig)
}
