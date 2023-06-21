package api_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/djmarrerajr/common-lib/services/api"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

const identityConfig = `
version: 1
allow-insecure-connections: true

security-checks:
  azuread-qa-aud-check:
    issuer-name: azuread-qa
    audience: spiffe://homedepot.dev/prop-tend-gc-lkup-srvr

paths:
  p-001:
    path: /health/*
    check-names:
      - azuread-qa-aud-check
`

type ApiTestSuite struct {
	suite.Suite

	server *api.Server
	addr   string
	cert   string
	key    string

	tmpdir string
	appctx shared.ApplicationContext
}

func (s *ApiTestSuite) SetupTest() {
	s.addr = "127.0.0.1"

	s.tmpdir = os.TempDir()
	s.cert = fmt.Sprintf("%s/temp.crt", s.tmpdir)
	s.key = fmt.Sprintf("%s/temp.key", s.tmpdir)

	s.appctx = shared.ApplicationContext{
		Logger: utils.NewLogger("INFO"),
	}
}

func (s *ApiTestSuite) TearDownTest() {
	os.RemoveAll(s.cert)
	os.RemoveAll(s.key)
}

func (s *ApiTestSuite) TestCreate_HttpServer() {
	var err error

	s.server, err = api.NewHttpServer(s.addr, "8080",
		api.WithTimeoutDurationSecs(1, 2, 3, 4),
	)

	s.NoError(err)
	s.Equal(fmt.Sprintf("%s:%s", s.addr, "8080"), s.server.Api.Addr)
	s.Equal(time.Second*1, s.server.Api.ReadTimeout)
	s.Equal(time.Second*2, s.server.Api.ReadHeaderTimeout)
	s.Equal(time.Second*3, s.server.Api.WriteTimeout)
	s.Equal(time.Second*4, s.server.Api.IdleTimeout)
}

func (s *ApiTestSuite) TestCreate_HttpsServer() {
	var err error

	err = s.genX509KeyPair()
	s.NoError(err)

	s.server, err = api.NewHttpsServer(s.addr, "8443", s.cert, s.key,
		api.WithMtlsEnforcedCaCert(s.cert),
	)

	s.NoError(err)
	s.Equal(fmt.Sprintf("%s:%s", s.addr, "8443"), s.server.Api.Addr)
	s.Equal(1, len(s.server.Api.TLSConfig.Certificates))
	s.NotNil(s.server.Api.TLSConfig.ClientCAs)
	s.Equal(tls.RequireAndVerifyClientCert, s.server.Api.TLSConfig.ClientAuth)
}

func (s *ApiTestSuite) TestCreate_HttpsServer_MissingKeyPair() {
	_, err := api.NewHttpsServer(s.addr, "8443", s.cert, s.key)

	s.Error(err)
	s.ErrorIs(err, os.ErrNotExist)
}

func (s *ApiTestSuite) TestConstructor_NewServerFromEnv_DefaultHttp() {
	var err error

	cleanup := s.setupEnviron(map[string]string{
		"THD_ID_CONFIG_DATA":      identityConfig,
		"THD_ID_CONFIG_DATA_TYPE": "yaml",
	})

	env := utils.GetEnviron()

	s.server, err = api.NewServerFromEnv(env, s.appctx)

	s.NoError(err)
	s.Equal(fmt.Sprintf("%s:%d", api.DefaultBindToAddress, api.DefaultHttpBindToPort), s.server.Api.Addr)
	s.Equal(api.DefaultReadTimeout, s.server.Api.ReadTimeout)
	s.Equal(api.DefaultReadHeaderTimeout, s.server.Api.ReadHeaderTimeout)
	s.Equal(api.DefaultWriteTimeout, s.server.Api.WriteTimeout)
	s.Equal(api.DefaultIdleTimeout, s.server.Api.IdleTimeout)
	s.NoError(s.checkRouteDefined("/health/live"))
	s.NoError(s.checkRouteDefined("/health/ready"))
	s.NoError(s.checkRouteDefined("/metrics"))

	s.T().Cleanup(cleanup)
}

func (s *ApiTestSuite) TestConstructor_NewServerFromEnv_DefaultHttps() {
	var err error

	err = s.genX509KeyPair()
	s.NoError(err)

	cleanup := s.setupEnviron(map[string]string{
		"THD_ID_CONFIG_DATA":      identityConfig,
		"THD_ID_CONFIG_DATA_TYPE": "yaml",
		api.ServerCertEnvKey:      s.cert,
		api.ServerKeyEnvKey:       s.key,
	})

	env := utils.GetEnviron()

	s.server, err = api.NewServerFromEnv(env, s.appctx)

	s.NoError(err)
	s.Equal(fmt.Sprintf("%s:%d", api.DefaultBindToAddress, api.DefaultHttpsBindToPort), s.server.Api.Addr)
	s.Equal(1, len(s.server.Api.TLSConfig.Certificates))
	s.Equal(api.DefaultReadTimeout, s.server.Api.ReadTimeout)
	s.Equal(api.DefaultReadHeaderTimeout, s.server.Api.ReadHeaderTimeout)
	s.Equal(api.DefaultWriteTimeout, s.server.Api.WriteTimeout)
	s.Equal(api.DefaultIdleTimeout, s.server.Api.IdleTimeout)
	s.NoError(s.checkRouteDefined("/health/live"))
	s.NoError(s.checkRouteDefined("/health/ready"))
	s.NoError(s.checkRouteDefined("/metrics"))

	s.T().Cleanup(cleanup)
}

func (s *ApiTestSuite) setupEnviron(envs map[string]string) (cleanup func()) {
	originalEnvs := map[string]string{}

	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvs[name] = originalValue
		}
		_ = os.Setenv(name, value)
	}

	return func() {
		for name := range envs {
			origValue, has := originalEnvs[name]
			if has {
				_ = os.Setenv(name, origValue)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}
}

// query the embedded router to see if the specified path is defined
func (s *ApiTestSuite) checkRouteDefined(path string) error {
	var currRoute *mux.Route

	// because it is possible to override a default route handler we need to check if it exists and
	// replace the handler because gorilla does not handle multiple route definitions well so we are
	// updating any route definitions that match our path so that they have the same handler...
	err := s.server.Api.Handler.(*mux.Router).Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
		return err
	}

	if currRoute == nil {
		return errors.New("route not defined")
	}

	return nil
}

func (s *ApiTestSuite) genX509KeyPair() error {
	os.RemoveAll(fmt.Sprintf("%s/temp.*", s.tmpdir))

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
	)

	template := x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * 10 * time.Hour),
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	crtBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return err
	}

	crtPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: crtBytes,
		},
	)

	err = os.WriteFile(s.cert, crtPEM, 0666)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.key, keyPEM, 0666)
	if err != nil {
		return err
	}

	return nil
}

func TestApi(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
