package api

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/gorilla/mux"

	"github.com/djmarrerajr/common-lib/utils"
)

// createNewTlsConfig is a convenience helper that allows us to
// create consistent, default, TLS Configurations
func createNewTlsConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
}

// WithLogger will assign the provided logger to the API
func WithLogger(logger utils.Logger) Option {
	return func(s *Server) {
		s.Logger = logger
	}
}

// WithTimeoutDurationSecs will update the appropriate Timeout duration where
// a positive integer value is provided
func WithTimeoutDurationSecs(read, readHeader, write, idle int) Option {
	return func(s *Server) {
		if read != 0 {
			s.Api.ReadTimeout = time.Duration(read) * time.Second
		}
		if readHeader != 0 {
			s.Api.ReadHeaderTimeout = time.Duration(readHeader) * time.Second
		}
		if write != 0 {
			s.Api.WriteTimeout = time.Duration(write) * time.Second
		}
		if idle != 0 {
			s.Api.IdleTimeout = time.Duration(idle) * time.Second
		}
	}
}

// WithTlsConfig will update the Server's TLS Configuration with
// the one that is provided
func WithTlsConfig(config *tls.Config) Option {
	return func(s *Server) {
		s.Api.TLSConfig = config
	}
}

// WithMtlsEnforcedCaCert will update the Server's TLS Configuration
// to ensure any clients that connect must provide a CA cert
func WithMtlsEnforcedCaCert(ca string) Option {
	return func(s *Server) {
		caCert, err := os.ReadFile(ca)
		if err != nil {
			s.Logger.Fatalf("unable to load ca cert: %v", err)
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caCert)

		if s.Api.TLSConfig == nil {
			s.Api.TLSConfig = createNewTlsConfig()
		}

		s.Api.TLSConfig.ClientCAs = certPool
		s.Api.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
}

// WithTlsConfig will update the Server's TLS Configuration
func WithCertificateAndKey(cert, key string) Option {
	return func(s *Server) {
		s.serverCert = cert
		s.serverKey = key

		cert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			s.Logger.Fatalf("loading x509 keypair: %v", err)
		}

		if s.Api.TLSConfig == nil {
			s.Api.TLSConfig = createNewTlsConfig()
		}

		s.Api.TLSConfig.Certificates = append(s.Api.TLSConfig.Certificates, cert)
	}
}

// WithRouteHandler will update the router/mux with a new path->handler
// ensuring that a route/mux exists prior to adding the mapping
//
// If no methods are specified we will default to GET
func WithRouteHandler(path string, handler http.HandlerFunc, methods ...string) Option {
	if len(methods) == 0 {
		methods = []string{http.MethodGet}
	}

	return func(s *Server) {
		if s.Api.Handler == nil {
			s.Api.Handler = mux.NewRouter()
		}

		defineOrReplaceRoute(s, path, handler, methods...)
	}
}

// WithRequestMiddleware will add a function to be called during the http.Request
// processing lifecycle... they are executed in the order in which they are added
func WithRequestMiddleware(fn mux.MiddlewareFunc) Option {
	return func(s *Server) {
		if s.Api.Handler == nil {
			s.Api.Handler = mux.NewRouter()
		}

		s.Logger.Debugf("adding new middleware function to the chain: %v", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
		s.Api.Handler.(*mux.Router).Use(fn)
	}
}

// WithPostShutdownCallback registers a function that will be invoked *after* the
// Server shutdown has been completed
func WithPostShutdownCallback(fn func()) Option {
	return func(s *Server) {
		s.Api.RegisterOnShutdown(fn)
	}
}

func defineOrReplaceRoute(s *Server, path string, handler http.HandlerFunc, methods ...string) {
	var currRoute *mux.Route

	// because it is possible to override a default route handler we need to check if it exists and
	// replace the handler because gorilla does not handle multiple route definitions well so we are
	// updating any route definitions that match our path so that they have the same handler...
	err := s.Api.Handler.(*mux.Router).Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		if tpl == path {
			currRoute = route
		}

		return nil
	})
	if err != nil {
		s.Logger.Fatalf("unable to define api route: %v", err)
	}

	if currRoute != nil {
		s.Logger.Debugf("redefining route for %s with %v", path, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		currRoute.Handler(handler)
	} else {
		s.Logger.Debugf("defining NEW route for %s with %v", path, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		s.Api.Handler.(*mux.Router).Handle(path, handler).Methods(methods...)
	}
}
