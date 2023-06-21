package db

import "time"

// nolint: unused
const (
	DefaultMaxIdleConnections = 3
	DefaultMaxOpenConnections = 10
	DefaultMaxConnLifeTime    = 10 * time.Second
	DefaultMaxConnIdleTime    = 5 * time.Second
)

// nolint: unused
const (
	CaCertEnvKey   = "DB_CA_CERT"
	UsernameEnvKey = "DB_USERNAME"
	UserCertEnvKey = "DB_USER_CERT"
	UserKeyEnvKey  = "DB_USER_KEY"

	DatabaseHostEnvKey = "DB_HOST_NAME"
	DatabasePortEnvKey = "DB_HOST_PORT"
	DatabaseNameEnvKey = "DB_DATABASE_NAME"

	DatabaseMaxIdleConnEnvKey = "DB_MAX_IDLE_CONNECTIONS"
	DatabaseMaxOpenConnEnvKey = "DB_MAX_OPEN_CONNECTIONS"
	DatabaseMaxIdleTimeEnvKey = "DB_MAX_IDLE_TIME_SECS"
	DatabaseMaxOpenTimeEnvKey = "DB_MAX_OPEN_TIME_SECS"
)
