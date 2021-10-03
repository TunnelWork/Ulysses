package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/TunnelWork/Ulysses/src/internal/conf"
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

func DBConnect(sconf conf.DatabaseConfig) (*sql.DB, error) {
	driverName := "mysql"
	// dsn = fmt.Sprintf("user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local", sconf.User, sconf.Passwd, sconf.Host, sconf.Port, sconf.Database)
	if mysqlAutoCommit {
		dsn += "&autocommit=true"
	}
	if sconf.CA != "" {
		dsn += "&tls=custom"
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(sconf.CA)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, ErrCannotAppendCert
		}
		if sconf.ClientKey != "" && sconf.ClientCert != "" {
			// Both Key and Cert are set. Go with customer cert.
			clientCert := make([]tls.Certificate, 0, 1)
			certs, err := tls.LoadX509KeyPair(sconf.ClientCert, sconf.ClientKey)
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
		} else if sconf.ClientKey == "" && sconf.ClientCert == "" {
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

func DBConnectWithContext(ctx context.Context, sconf conf.DatabaseConfig) (*sql.DB, error) {
	var dbConn *sql.DB
	var err error

	dbDone := make(chan bool)

	go func() {
		dbConn, err = DBConnect(sconf)
		dbDone <- true
	}()

	select {
	case <-dbDone:
		return dbConn, err
	case <-ctx.Done():
		return dbConn, ctx.Err()
	}
}
