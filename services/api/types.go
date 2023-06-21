package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

var _ shared.Servable = new(Server)

type Server struct {
	AppCtx shared.ApplicationContext

	Api    *http.Server
	Logger utils.Logger

	serverCert string
	serverKey  string
}

func (s Server) Start(ctx context.Context, grp *errgroup.Group) error {
	if s.Api.TLSConfig != nil {
		grp.Go(func() error {
			s.Logger.Infof("listening for tls connections on %s", s.Api.Addr)
			if err := s.Api.ListenAndServeTLS(s.serverCert, s.serverKey); err != http.ErrServerClosed {
				s.Logger.Fatalf("unable to shutdown server: %v", err)
				return err
			}

			s.Logger.Infof("tls server no longer listening on %s", s.Api.Addr)
			return nil
		})

	} else {

		grp.Go(func() error {
			s.Logger.Infof("listening for connections on %s", s.Api.Addr)
			if err := s.Api.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("unable to shutdown server: %v", err)
				return err
			}

			s.Logger.Infof("server no longer listening on %s", s.Api.Addr)
			return nil
		})
	}

	return nil
}

func (s Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	s.Logger.Infof("waiting %s for existing connections to terminate", DefaultShutdownTimeout)
	if err := s.Api.Shutdown(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.Logger.Infof("shutdown timeout exceeded - remaining connections terminated")
			return nil
		}

		return err
	}

	s.Logger.Infof("api termination complete")
	return nil
}

func (s Server) DefineRoute(path string, handler http.HandlerFunc, methods ...string) {
	defineOrReplaceRoute(&s, path, handler, methods...)
}

func (s Server) DefineRequestHandler(path string, handler shared.RequestHandlerFunc, reqStruct any, methods ...string) {
	ctxHandler := ContextualHandler{&s.AppCtx, handler, reqStruct}

	defineOrReplaceRoute(&s, path, ctxHandler.ServeHTTP, methods...)
}
