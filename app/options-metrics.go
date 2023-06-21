package app

import (
	"github.com/djmarrerajr/common-lib/observability/metrics"
	"github.com/djmarrerajr/common-lib/utils"
)

func WithMetricsCollector(Collector metrics.Collector) Option {
	return func(a *application) {
		a.AppContext.Collector = Collector
	}
}

func WithMetricsCollectorFromEnv(env utils.Environ, options ...metrics.Option) Option {
	return func(a *application) {
		e, err := metrics.NewCollectorFromEnv(env, a.AppContext.Logger, options...)
		if err != nil {
			a.AppContext.Logger.Fatalf("unable to create metrics Collector:  %v", err)
		}

		a.AppContext.Collector = e
	}
}
