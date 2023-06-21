package api

import "time"

// nolint: unused
const (
	DefaultReadTimeout       = 15 * time.Second
	DefaultReadHeaderTimeout = 15 * time.Second
	DefaultWriteTimeout      = 15 * time.Second
	DefaultIdleTimeout       = 15 * time.Second
	DefaultShutdownTimeout   = 15 * time.Second

	DefaultBindToAddress   = "0.0.0.0"
	DefaultHttpBindToPort  = 8080
	DefaultHttpsBindToPort = 8443
)

// nolint: unused
const (
	BindToAddressEnvKey = "API_BIND_ADDRESS"
	BindToPortEnvKey    = "API_BIND_PORT"
	ServerCertEnvKey    = "API_SERVER_CERT"
	ServerKeyEnvKey     = "API_SERVER_KEY"

	ReadTimeoutEnvKey       = "API_READ_TIMEOUT"
	ReadHeaderTimeoutEnvKey = "API_READHEADER_TIMEOUT"
	WriteTimeoutEnvKey      = "API_WRITE_TIMEOUT"
	IdleTimeoutEnvKey       = "API_IDLE_TIMEOUT"
)

// nolint: unused
const (
	HeaderContentType   = "Content-Type"
	HeaderContentLength = "Content-Length"
	HeaderRequestId     = "X-Request-Id"
)

// nolint: unused
const (
	ValueTextPlain       = "text/plain"
	ValueApplicationJson = "application/json"
)
