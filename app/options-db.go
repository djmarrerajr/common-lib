package app

import (
	"github.com/djmarrerajr/common-lib/services/db"
	"github.com/djmarrerajr/common-lib/services/db/cockroach"
	"github.com/djmarrerajr/common-lib/utils"
)

func WithCockroachDB(database db.Adapter) Option {
	return func(a *application) {
		a.AppContext.Database = database
	}
}

func WithCockroachDBFromEnv(env utils.Environ, options ...cockroach.Option) Option {
	return func(a *application) {
		s, err := cockroach.NewAdapterFromEnv(env, *a.AppContext, options...)
		if err != nil {
			a.AppContext.Logger.Fatalf("unable to create db adapter:  %v", err)
		}

		a.AppContext.Database = s
	}
}
