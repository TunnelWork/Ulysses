package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/go-sql-driver/mysql"
)

const (
	mysqlAutoCommit = true
)

func DBConnected(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		db.Close()
		return false
	}
	return true
}

func DBConnect(sconf MysqlConf) (*sql.DB, error) {
	driverName := "mysql"
	// dsn = fmt.Sprintf("user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local", sconf.mysqlUser, sconf.mysqlPasswd, sconf.mysqlHost, sconf.mysqlPort, sconf.mysqlDatabase)
	if mysqlAutoCommit {
		dsn += "&autocommit=true"
	}
	if sconf.mysqlCAPath != "" {
		dsn += "&tls=custom"
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(sconf.mysqlCAPath)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, ErrCannotAppendCert
		}
		if sconf.mysqlKeyPath != "" && sconf.mysqlCertPath != "" {
			// Both Key and Cert are set. Go with customer cert.
			clientCert := make([]tls.Certificate, 0, 1)
			certs, err := tls.LoadX509KeyPair(sconf.mysqlCertPath, sconf.mysqlKeyPath)
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
		} else if sconf.mysqlKeyPath == "" && sconf.mysqlCertPath == "" {
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

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	if connected := DBConnected(db); !connected {
		return nil, ErrMySQLNoConn
	}

	return db, nil
}
