package uconf

import (
	"database/sql"
	"encoding/json"
)

type HttpServerConfig struct {
	// HTTP Server
	HTTPHost    string `json:"http_host"`    // Which IP to listen on. e.g., 127.0.0.1
	HTTPPort    uint16 `json:"http_port"`    // Which port to listen on. e.g., 8080
	URLDomain   string `json:"url_domain"`   // e.g., ulysses.tunnel.work
	URLPrefix   string `json:"url_prefix"`   // URL Prefix for API. e.g., /api
	URLComplete string `json:"url_complete"` // Complete URL for your publicly accessible Web Client, including protocol and port. e.g., https://ulysses.tunnel.work:8443
}

// This is solely for debugging purpose.
var DefaultHttpServerConfig = HttpServerConfig{
	HTTPHost:    "127.0.0.1",
	HTTPPort:    9891,
	URLDomain:   "ulysses-integration.tunnel.work",
	URLPrefix:   "/api",
	URLComplete: "https://ulysses-integration.tunnel.work", // You are advised to shield Ulysses instance under a web proxy server, such as nginx.
}

func (c HttpServerConfig) String() string {
	// Marshal to JSON
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(b)
}

func LoadHttpServerConfig() (HttpServerConfig, error) {
	stmt, err := sqlStatement(`SELECT config_content FROM dbprefix_config WHERE config_name = ?`)
	if err != nil {
		return DefaultHttpServerConfig, err
	}
	defer stmt.Close()

	var config_content string
	err = stmt.QueryRow("http").Scan(&config_content)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert default config
			stmt, err = sqlStatement(`INSERT INTO dbprefix_config (config_name, config_content) VALUES (?, ?)`)
			if err != nil {
				return DefaultHttpServerConfig, err
			}
			defer stmt.Close()

			_, err = stmt.Exec("http", DefaultHttpServerConfig.String())
			if err != nil {
				return DefaultHttpServerConfig, err
			}
			return DefaultHttpServerConfig, nil
		} else {
			return DefaultHttpServerConfig, err
		}
	}

	// unmarshal
	var c HttpServerConfig
	err = json.Unmarshal([]byte(config_content), &c)
	if err != nil {
		return DefaultHttpServerConfig, err
	}
	return c, nil
}
