package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/djmarrerajr/common-lib/shared"
)

// MetricsMiddleware will integrate with the metrics collector service to create
// and increment a standard set of obversability metrics for each, registered,
// api endpoint
func MetricsMiddleware(appCtx shared.ApplicationContext) mux.MiddlewareFunc {
	collector := appCtx.Collector

	requestByPath := collector.NewDimensionedCounter("requests_total", "path")
	responseStatusByPath := collector.NewDimensionedCounter("response_status", "path", "statusCode")
	responseErrorsByPath := collector.NewDimensionedCounter("response_errors", "path", "errorType")
	responseTimeByPath := collector.NewDimensionedGauge("response_time", "path")

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.ReplaceAll(r.URL.Path[1:], "/", "_")

			mw := &metricsResponseWriter{writer: w}
			st := time.Now()

			requestByPath.WithLabelValues(path).Inc()

			h.ServeHTTP(mw, r)

			if mw.errorType != "" {
				responseErrorsByPath.WithLabelValues(path, string(mw.errorType)).Inc()
			}

			responseStatusByPath.WithLabelValues(path, fmt.Sprint(mw.code)).Inc()
			responseTimeByPath.WithLabelValues(path).Set(float64(time.Since(st).Milliseconds()))
		})
	}
}
