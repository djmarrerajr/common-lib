package api

// AddRequestInfoToContext will extract key bits of information from the
// request (i.e. a unique request-id, etc.) and add them to the context
// so that they can be logged
// func AddRequestInfoToContext(appCtx shared.ApplicationContext) mux.MiddlewareFunc {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			reqId := r.Header.Get(HeaderRequestId)
// 			if reqId == "" {
// 				reqId = uuid.New().String()
// 			}

// 			appCtx.RootCtx = utils.AddMapToContext(appCtx.RootCtx, utils.FieldMap{
// 				"request-url": r.URL.Path,
// 				"request-id":  reqId,
// 			})

// 			appCtx.Logger.WithCtx(appCtx.RootCtx).Debugf("request received")
// 			h.ServeHTTP(w, r)
// 			appCtx.Logger.WithCtx(appCtx.RootCtx).Debugf("request completed")
// 		})
// 	}
// }

// MetricsMiddleware will integrate with the metrics collector service to create
// and increment a standard set of obversability metrics for each, registered,
// apo endpoint
// func MetricsMiddleware(collector metrics.Collector) mux.MiddlewareFunc {
// 	requestByPath := collector.NewDimensionedCounter("requests_total", "path")
// 	responseStatusByPath := collector.NewDimensionedCounter("response_status", "path", "statusCode")
// 	responseTimeByPath := collector.NewDimensionedGauge("response_time", "path")

// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			path := strings.ReplaceAll(r.URL.Path[1:], "/", "_")

// 			mw := &metricsResponseWriter{writer: w}
// 			st := time.Now()

// 			requestByPath.WithLabelValues(path).Inc()

// 			h.ServeHTTP(mw, r)

// 			responseStatusByPath.WithLabelValues(path, fmt.Sprint(mw.code)).Inc()
// 			responseTimeByPath.WithLabelValues(path).Set(float64(time.Since(st).Milliseconds()))
// 		})
// 	}
// }