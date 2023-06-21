package app

import (
	"net/http"

	"github.com/djmarrerajr/common-lib/services/api"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

func WithApiServer(s *api.Server) Option {
	return func(a *application) {
		a.AppContext.Server = s
		a.AppContext.Server.(*api.Server).Logger = a.AppContext.Logger.Named("api")
	}
}

func WithApiServerFromEnv(env utils.Environ, options ...api.Option) Option {
	return func(a *application) {
		s, err := api.NewServerFromEnv(env, *a.AppContext, options...)
		if err != nil {
			a.AppContext.Logger.Fatalf("unable to create api server:  %v", err)
		}

		a.AppContext.Server = s
	}
}

// WithRouteHandler is a convenience function that allows for the definition of an API
// route handler without having to replace the default API
func WithRouteHandler(path string, handler http.HandlerFunc, methods ...string) Option {
	return func(a *application) {
		a.AppContext.Server.DefineRoute(path, handler, methods...)
	}
}

// WithRequestHandler is a convenience function that allows for the definition of an API
// request handler without having to replace the default API
func WithRequestHandler(path string, handler shared.RequestHandlerFunc, reqStruct any, methods ...string) Option {
	return func(a *application) {
		a.AppContext.Server.DefineRequestHandler(path, handler, reqStruct, methods...)
	}
}

// WithEndpoint is a convenience function that allows for the definition of an API
// route handler without having to replace the default API
// func WithEndpoint(path string, reqType any, handler shared.HandlerFunc, methods ...string) Option {
// 	return func(a *application) {
// 		a.AppContext.Server.DefineRoute(*a.AppContext, path, reqType, handler, methods...)
// 	}
// }
