package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

// MetricsMiddleware will integrate with the metrics collector service to create
// and increment a standard set of obversability metrics for each, registered,
// api endpoint
func MetricsMiddleware(appCtx shared.ApplicationContext) mux.MiddlewareFunc {
	collector := appCtx.Collector

	filterLabels := []string{shared.EnvironContextKey, shared.HostnameContextKey, shared.AppNameContextKey, shared.AppVersionContextKey}

	// define out standard set of api metrics...
	requestByPath := collector.NewDimensionedCounter("requests_total", append(filterLabels, "path")...)
	responseStatusByPath := collector.NewDimensionedCounter("response_status", append(filterLabels, "path", "statusCode")...)
	responseErrorsByPath := collector.NewDimensionedCounter("response_errors", append(filterLabels, "path", "errorType")...)
	responseTimeByPath := collector.NewDimensionedGauge("response_time", append(filterLabels, "path")...)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.ReplaceAll(r.URL.Path[1:], "/", "_")

			mw := &metricsResponseWriter{writer: w}
			st := time.Now()

			h.ServeHTTP(mw, r)

			// grab values for our standard metric labels...
			envn, _ := utils.GetFieldValueFromContext[string](appCtx.RootCtx, shared.EnvironContextKey)
			host, _ := utils.GetFieldValueFromContext[string](appCtx.RootCtx, shared.HostnameContextKey)
			appl, _ := utils.GetFieldValueFromContext[string](appCtx.RootCtx, shared.AppNameContextKey)
			vrsn, _ := utils.GetFieldValueFromContext[string](appCtx.RootCtx, shared.AppVersionContextKey)

			filterValues := []string{envn, host, appl, vrsn}

			// now increment our standard metrics...
			requestByPath.WithLabelValues(append(filterValues, path)...).Inc()
			responseStatusByPath.WithLabelValues(append(filterValues, path, fmt.Sprint(mw.code))...).Inc()
			responseTimeByPath.WithLabelValues(append(filterValues, path)...).Set(float64(time.Since(st).Milliseconds()))
			if mw.errorType != "" {
				responseErrorsByPath.WithLabelValues(append(filterValues, path, string(mw.errorType))...).Inc()
			}
		})
	}
}
