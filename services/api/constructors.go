package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/djmarrerajr/common-lib/errs"
	"github.com/djmarrerajr/common-lib/observability/tracing"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

// NewHttpServer will create and return a configurable Server that
// wraps an underlying non-TLS enabled http.Server
func NewHttpServer(addr, port string, options ...Option) (*Server, error) {
	server := createServer(addr, port, options...)

	return server, nil
}

// NewHttpsServer will create and return a configurable Server that
// wraps an underlying TLS enabled http.Server
func NewHttpsServer(addr, port, cert, key string, options ...Option) (*Server, error) {
	if _, err := os.Stat(cert); err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	if _, err := os.Stat(key); err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	options = append([]Option{WithCertificateAndKey(cert, key)}, options...)

	server := createServer(addr, port, options...)

	return server, nil
}

// NewServerFromEnv will instantiate and return a new Server that has been configured using
// the values retrieved from the environment
//
// It will, by default, have the following endpoints available:
// ... /health	- returns an HTTP-200
// ... /metrics	- returns the full set of Prometheus metrics being collected
func NewServerFromEnv(env utils.Environ, appCtx shared.ApplicationContext, options ...Option) (*Server, error) {
	logger := appCtx.Logger.Named("api")
	newopt := []Option{WithLogger(logger)}

	// get the paths to our server cert/key if available
	cert, _ := env.Get(ServerCertEnvKey)
	if _, err := os.Stat(cert); cert != "" && err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	key, _ := env.Get(ServerKeyEnvKey)
	if _, err := os.Stat(key); key != "" && err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	// pull the address to which we should bind from the env
	addr, OK := env.Get(BindToAddressEnvKey)
	if !OK {
		addr = DefaultBindToAddress
	}

	// pull the port to which we should bind from the env
	port, OK, err := env.GetInt(BindToPortEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	} else if !OK {
		port = DefaultHttpBindToPort

		if cert != "" && key != "" {
			port = DefaultHttpsBindToPort
			newopt = append(newopt, WithCertificateAndKey(cert, key))
		}
	}

	// grab any timeout values from the env
	readTimeout, _, err := env.GetInt(ReadTimeoutEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	readHeaderTimeout, _, err := env.GetInt(ReadHeaderTimeoutEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	writeTimeout, _, err := env.GetInt(WriteTimeoutEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	idleTimeout, _, err := env.GetInt(IdleTimeoutEnvKey)
	if err != nil {
		return nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	newopt = append(newopt,
		WithTimeoutDurationSecs(readTimeout, readHeaderTimeout, writeTimeout, idleTimeout),
		WithRequestMiddleware(tracing.RequestTracing(appCtx, "/health", "/metrics")),
		WithRequestMiddleware(MetricsMiddleware(appCtx)),
		WithRouteHandler("/health", defaultHealthCheckHandler),
		WithRouteHandler("/metrics", http.HandlerFunc(promhttp.Handler().ServeHTTP)),
	)

	newopt = append(newopt, options...)

	server := createServer(addr, fmt.Sprint(port), newopt...)
	server.AppCtx = appCtx

	return server, nil
}

// createServer will create and return a configurable Server that
// encapsulates an underlying http.Server with all of the provided
// options having been applied
func createServer(addr, port string, options ...Option) *Server {
	server := Server{
		Api: &http.Server{
			Addr:              fmt.Sprintf("%s:%s", addr, port),
			ReadTimeout:       DefaultReadTimeout,
			ReadHeaderTimeout: DefaultReadHeaderTimeout,
			WriteTimeout:      DefaultWriteTimeout,
			IdleTimeout:       DefaultIdleTimeout,
		},
	}

	for _, option := range options {
		option(&server)
	}

	return &server
}
