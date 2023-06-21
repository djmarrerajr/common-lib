package cockroach

import (
	"time"

	"github.com/djmarrerajr/common-lib/utils"
)

type Option func(*CockroachDB)

func WithLogger(logger utils.Logger) Option {
	return func(cd *CockroachDB) {
		cd.logger = logger
	}
}

func WithConnectionInfo(host, port, user string) Option {
	return func(cd *CockroachDB) {
		cd.host = host
		cd.port = port
		cd.user = user
	}
}

func WithConnectionLimits(maxConn, idleConn, maxTime, idleTime int) Option {
	return func(cd *CockroachDB) {
		if maxConn != 0 {
			cd.maxConn = maxConn
		}
		if idleConn != 0 {
			cd.idleConn = idleConn
		}
		if maxTime != 0 {
			cd.maxTime = time.Duration(maxTime) * time.Second
		}
		if idleTime != 0 {
			cd.idleTime = time.Duration(idleTime) * time.Second
		}
	}
}

func WithCertificateInfo(ca, cert, key string) Option {
	return func(cd *CockroachDB) {
		cd.ca = ca
		cd.cert = cert
		cd.key = key
	}
}
