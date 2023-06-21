package cockroach

import (
	"fmt"
	"os"

	"github.com/djmarrerajr/common-lib/errs"
	"github.com/djmarrerajr/common-lib/services/db"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

func NewAdapterFromEnv(env utils.Environ, appCtx shared.ApplicationContext, options ...Option) (*CockroachDB, error) {
	logger := appCtx.Logger.Named("db")
	newopt := []Option{WithLogger(logger.WithCtx(appCtx.RootCtx))}

	// get the paths to our server cert/key if available
	cacert, err := env.GetRequired(db.CaCertEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	} else {
		if _, err := os.Stat(cacert); cacert != "" && err != nil {
			return nil, errs.WithType(err, errs.ErrTypeConfiguration)
		}
	}

	cert, err := env.GetRequired(db.UserCertEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	} else {
		if _, err := os.Stat(cert); cert != "" && err != nil {
			return nil, errs.WithType(err, errs.ErrTypeConfiguration)
		}
	}

	key, err := env.GetRequired(db.UserKeyEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	} else {
		if _, err := os.Stat(key); key != "" && err != nil {
			return nil, errs.WithType(err, errs.ErrTypeConfiguration)
		}
	}

	// get the connection details from the env
	host, err := env.GetRequired(db.DatabaseHostEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)

	}

	port, _, err := env.GetInt(db.DatabasePortEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	dbuser, err := env.GetRequired(db.UsernameEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)

	}

	dbname, err := env.GetRequired(db.DatabaseNameEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)

	}

	// get the connection and time limits
	idleConn, OK, err := env.GetInt(db.DatabaseMaxIdleConnEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	} else if !OK {
		idleConn = db.DefaultMaxIdleConnections
	}

	idleTime, _, err := env.GetInt(db.DatabaseMaxIdleTimeEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	maxConn, OK, err := env.GetInt(db.DatabaseMaxOpenConnEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	} else if !OK {
		maxConn = db.DefaultMaxOpenConnections
	}

	maxTime, _, err := env.GetInt(db.DatabaseMaxOpenTimeEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	newopt = append(newopt,
		WithConnectionInfo(host, fmt.Sprint(port), dbuser),
		WithConnectionLimits(maxConn, idleConn, maxTime, idleTime),
		WithCertificateInfo(cacert, cert, key),
	)

	return NewCockroachDB(dbname, newopt...), nil
}

func NewCockroachDB(database string, options ...Option) *CockroachDB {
	db := CockroachDB{
		database: database,
	}

	for _, option := range options {
		option(&db)
	}

	return &db
}
