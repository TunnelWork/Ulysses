package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-sql-driver/mysql"
)

// MysqlConnect is a STATELESS struct used to generate new MySQL connections
type MysqlConnector struct {
	conf    DatabaseConfig
	timeout time.Duration // !unimplemented
}

// NewMysqlConnector returns a valid pointer to a
// MysqlConnector struct when conf is valid (able to
// establish mysql connections)
func NewMysqlConnector(conf DatabaseConfig) *MysqlConnector {
	mysqlConnector := MysqlConnector{
		conf: conf,
	}

	conn, err := mysqlConnector.Conn()
	if err != nil || conn.Ping() != nil {
		return nil
	}

	return &mysqlConnector
}

// Conn() creates the *sql.DB using the DatabaseConfig stored in
// current MysqlConnector
func (mc *MysqlConnector) Conn() (*sql.DB, error) {
	driverName := "mysql"
	// dsn = fmt.Sprintf("user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local", mc.conf.User, mc.conf.Passwd, mc.conf.Host, mc.conf.Port, mc.conf.Database)
	if mc.conf.MysqlAutoCommit {
		dsn += "&autocommit=true"
	}
	if mc.conf.CA != "" {
		dsn += "&tls=custom"
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(mc.conf.CA)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, ErrCannotAppendCert
		}
		if mc.conf.ClientKey != "" && mc.conf.ClientCert != "" {
			// Both Key and Cert are set. Go with customer cert.
			clientCert := make([]tls.Certificate, 0, 1)
			certs, err := tls.LoadX509KeyPair(mc.conf.ClientCert, mc.conf.ClientKey)
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
		} else if mc.conf.ClientKey == "" && mc.conf.ClientCert == "" {
			// Neither Key or Cert is set. Proceed without customer cert.
			mysql.RegisterTLSConfig("custom", &tls.Config{
				// ServerName: "example.com",
				RootCAs:    rootCertPool,
				MinVersion: tls.VersionTLS12,
				MaxVersion: 0,
			})
		} else {
			// one of Key or Cert is set but not both, which is ILLEGAL.
			return nil, ErrIncompleteConf
		}
	}

	var db *sql.DB
	var err error
	db, err = sql.Open(driverName, dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
