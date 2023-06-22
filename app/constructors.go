package app

import (
	"context"
	"os"

	"github.com/djmarrerajr/common-lib/errs"
	"github.com/djmarrerajr/common-lib/observability/metrics"
	"github.com/djmarrerajr/common-lib/observability/traces"
	"github.com/djmarrerajr/common-lib/services/api"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
	"github.com/go-playground/validator"
)

// NewWithApiFromEnv will instantiate and return a standardized, albeit
// functionally limited, application that includes:
//   - a structured logger
//   - a signal handler (USR1 = toggle debug logging, INT = shutdown)
//   - a basic HTTP API that:
//     ... responds to '/health' with an HTTP-200
//     ... responds to '/metrics' with Prometheus data
//
// NOTE: A Database adapter can be added to the base application via the
// app.Option (i.e. WithCockroachDBFromEnv)
func NewWithApiFromEnv(env utils.Environ, opts ...Option) (*application, error) {
	var err error

	app, err := createInitialApplication(env)
	if err != nil {
		return nil, err
	}

	app.AppContext.Collector, err = metrics.NewCollectorFromEnv(env, app.name)
	if err != nil {
		return nil, errs.Wrap(err, errs.ErrTypeConfiguration, "while instantiating metrics collector")
	}

	app.AppContext.Tracer, app.AppContext.Closer, err = traces.NewTracerFromEnv(env, *app.AppContext, app.name, app.version)
	if err != nil {
		return nil, errs.Wrap(err, errs.ErrTypeConfiguration, "while instantiating tracer")
	}

	app.AppContext.Server, err = api.NewServerFromEnv(env, *app.AppContext)
	if err != nil {
		return nil, errs.Wrap(err, errs.ErrTypeConfiguration, "while instantiating api")
	}

	for _, opt := range opts {
		opt(&app)
	}

	return &app, nil
}

// createInitialApplication will return an instance of our application
// struct that has been initialized with the basic elements we need in
// order to extend its functionality
func createInitialApplication(env utils.Environ) (application, error) {
	appName, err := env.GetRequired(AppNameEnvKey)
	if err != nil {
		return application{}, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	appVrsn, err := env.GetRequired(AppVersionEnvKey)
	if err != nil {
		return application{}, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	ctx := context.Background()
	ctx = utils.AddMapToContext(ctx, utils.FieldMap{
		"appName":    appName,
		"appVersion": appVrsn,
	})

	commit, OK := env.Get(AppCommitEnvKey)
	if OK {
		ctx = utils.AddFieldToContext(ctx, "commitId", commit)
	}

	return application{
		name:    appName,
		version: appVrsn,
		commit:  commit,

		env:            env,
		signalHandlers: make(map[os.Signal]signalHandler),

		AppContext: &shared.ApplicationContext{
			RootCtx:   ctx,
			Logger:    utils.NewLoggerFromEnv().Named(appName).WithCtx(ctx),
			Validator: validator.New(),
		},
	}, nil
}
