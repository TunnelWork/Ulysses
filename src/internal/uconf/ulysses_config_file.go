package uconf

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var (
	db        *sql.DB
	tblPrefix string
)

// UlyssesConfigFile is equivalent to a ulysses.yml file
type UlyssesConfigFile struct {
	/**************** MySQL ****************/
	// Mandatory
	Host     string `yaml:"host"` // For IPv6, use brackets enclosed representation. e.g., [::]
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Passwd   string `yaml:"passwd"`
	// Optional
	CA         string `yaml:"ca_cert"`     // If unset, ClientKey/ClientCert must not be set.
	ClientKey  string `yaml:"client_key"`  // If unset, ClientCert must not be set.
	ClientCert string `yaml:"client_cert"` // If unset, ClientKey must not be set.
	TblPrefix  string `default:"ulysses_" yaml:"table_prefix"`

	/**************** Security ****************/
	SecuritySeed string `yaml:"security_seed"` // Must be long enough to generate a secure encryption key.
}

func LoadConfigFromFile(configPath string) (UlyssesConfigFile, error) {
	var uc UlyssesConfigFile
	err := defaults.Set(&uc)
	if err != nil {
		return uc, err
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return uc, err
	}

	err = yaml.Unmarshal(content, &uc)
	return uc, err
}

func sqlStatement(query string) (*sql.Stmt, error) {
	prefixUpdatedQuery := strings.ReplaceAll(query, "dbprefix_", tblPrefix)
	return db.Prepare(prefixUpdatedQuery)
}

// DB() creates the *sql.DB to Ulysses Main Database
func (ucf UlyssesConfigFile) db() (*sql.DB, error) {
	driverName := "mysql"
	// dsn = fmt.Sprintf("user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&autocommit=true", ucf.User, ucf.Passwd, ucf.Host, ucf.Port, ucf.Database)

	if ucf.CA != "" {
		dsn += "&tls=custom"
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(ucf.CA)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, errors.New("uconf: AppendCertsFromPEM() has failed")
		}
		if ucf.ClientKey != "" && ucf.ClientCert != "" {
			// Both Key and Cert are set. Go with customer cert.
			clientCert := make([]tls.Certificate, 0, 1)
			certs, err := tls.LoadX509KeyPair(ucf.ClientCert, ucf.ClientKey)
			if err != nil {
				return nil, err
			}
			clientCert = append(clientCert, certs)
			mysql.RegisterTLSConfig("custom", &tls.Config{
				// ServerName: "example.com",
				RootCAs:      rootCertPool,
				Certificates: clientCert,
				MinVersion:   tls.VersionTLS12,
				MaxVersion:   0,
			})
		} else if ucf.ClientKey == "" && ucf.ClientCert == "" {
			// Neither Key or Cert is set. Proceed without customer cert.
			mysql.RegisterTLSConfig("custom", &tls.Config{
				// ServerName: "example.com",
				RootCAs:    rootCertPool,
				MinVersion: tls.VersionTLS12,
				MaxVersion: 0,
			})
		} else {
			// one of Key or Cert is set but not both, which is ILLEGAL.
			return nil, errors.New("uconf: must set both client_key and client_cert or set neither")
		}
	}

	var db *sql.DB
	var err error
	db, err = sql.Open(driverName, dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}

func (ucf UlyssesConfigFile) LoadCompleteConfig() (CompleteConfig, error) {
	var err error
	db, err = ucf.db()
	if err != nil {
		return CompleteConfig{}, err
	}

	// Make sure table is created.
	stmt, err := sqlStatement(`CREATE TABLE IF NOT EXISTS dbprefix_config (
		id INT NOT NULL AUTO_INCREMENT,
		config_name VARCHAR(32) NOT NULL,
		config_content TEXT NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY (config_name)
	)`)
	if err != nil {
		return CompleteConfig{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return CompleteConfig{}, err
	}

	// Load all configs.

	http, err := LoadHttpServerConfig()
	if err != nil {
		return CompleteConfig{}, err
	}

	logger, err := LoadLoggerConfig()
	if err != nil {
		return CompleteConfig{}, err
	}

	cc := CompleteConfig{
		Mysql: MysqlClientConfig{
			host:     ucf.Host,
			port:     ucf.Port,
			database: ucf.Database,
			user:     ucf.User,
			passwd:   ucf.Passwd,

			ca:         ucf.CA,
			clientKey:  ucf.ClientKey,
			clientCert: ucf.ClientCert,
			TblPrefix:  ucf.TblPrefix,
		},
		Http:   http,
		Logger: logger,
		Security: SecurityModuleConfig{
			secSeed: ucf.SecuritySeed,
		},
	}

	return cc, nil
}
