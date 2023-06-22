package metrics

import (
	"strings"

	"github.com/djmarrerajr/common-lib/utils"
)

func NewCollectorFromEnv(env utils.Environ, appName string, options ...Option) (*PrometheusCollector, error) {
	Collector := &PrometheusCollector{
		appName:  strings.ReplaceAll(appName, "-", "_"),
		counters: make(map[string]any),
	}

	return Collector, nil
}
