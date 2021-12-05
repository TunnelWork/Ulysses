package uconf

import (
	"database/sql"
	"encoding/json"

	"github.com/TunnelWork/Ulysses.Lib/logging"
)

type LoggerConfig logging.LoggerConfig

var DefaultLoggerConfig = LoggerConfig{
	Verbose:  true,
	Filepath: "./log/ulysses_debug.log",
	Level:    logging.LvlInfo,
}

func (l LoggerConfig) String() string {
	// marshal to json
	b, err := json.Marshal(l)
	if err != nil {
		return ""
	}
	return string(b)
}

func LoadLoggerConfig() (logging.LoggerConfig, error) {
	stmt, err := sqlStatement(`SELECT config_content FROM dbprefix_config WHERE config_name = ?`)
	if err != nil {
		return logging.LoggerConfig(DefaultLoggerConfig), err
	}
	defer stmt.Close()

	var config_content string
	err = stmt.QueryRow(`logger`).Scan(&config_content)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert default config
			stmt, err = sqlStatement(`INSERT INTO dbprefix_config (config_name, config_content) VALUES (?, ?)`)
			if err != nil {
				return logging.LoggerConfig(DefaultLoggerConfig), err
			}
			defer stmt.Close()

			_, err = stmt.Exec("logger", DefaultLoggerConfig.String())
			if err != nil {
				return logging.LoggerConfig(DefaultLoggerConfig), err
			}
			return logging.LoggerConfig(DefaultLoggerConfig), nil
		} else {
			return logging.LoggerConfig(DefaultLoggerConfig), err
		}
	}

	// unmarshal
	var config LoggerConfig
	err = json.Unmarshal([]byte(config_content), &config)
	if err != nil {
		return logging.LoggerConfig(DefaultLoggerConfig), err
	}
	return logging.LoggerConfig(config), nil
}
