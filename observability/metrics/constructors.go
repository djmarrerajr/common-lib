package metrics

import (
	"github.com/djmarrerajr/common-lib/utils"
)

func NewCollectorFromEnv(env utils.Environ, logger utils.Logger, options ...Option) (*PrometheusCollector, error) {
	Collector := &PrometheusCollector{
		counters: make(map[string]any),
	}

	return Collector, nil
}
