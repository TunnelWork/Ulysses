package main

import (
	"io/ioutil"

	"github.com/TunnelWork/Ulysses/src/internal/db"
)

func initDB() {
	content, err := ioutil.ReadFile(dbConfigPath)
	if err != nil {
		Fatal("initDB(): ", err)
	}

	dbConfig = db.LoadMySqlConf(content)

	// DB Conn Livess Test
	_, err = db.DBConnect(dbConfig)
	if err != nil {
		Fatal("initDB(): ", err)
	}
	Debug("initDB(): PASS")
}
