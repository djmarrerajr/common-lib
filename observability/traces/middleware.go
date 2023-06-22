package traces

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

func RequestTracing(appCtx shared.ApplicationContext, routesToSuppress ...string) mux.MiddlewareFunc {
	routes := make(map[string]struct{})

	for _, r := range routesToSuppress {
		routes[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			endpoint := r.URL.String()
			reqID := r.Header.Get(HeaderRequestId)
			if reqID == "" {
				reqID = uuid.NewString()
			}

			r = r.WithContext(utils.AddFieldToContext(r.Context(), "requestId", reqID))

			if _, exists := routes[endpoint]; !exists {
				var span opentracing.Span

				tracer := opentracing.GlobalTracer()
				// begin := time.Now()

				// If the incoming request is carrying any opentracing context information
				// extract it so it can be used...
				carrier := opentracing.HTTPHeadersCarrier(r.Header)
				clientContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)

				// If the request carried any context, the span will also carry this context
				// ... otherwise, we create a root span.
				if err == nil {
					span = tracer.StartSpan(endpoint, ext.RPCServerOption(clientContext))
				} else {
					span = tracer.StartSpan(endpoint)
				}

				ext.Component.Set(span, r.URL.Scheme)
				ext.HTTPMethod.Set(span, r.Method)
				ext.HTTPUrl.Set(span, endpoint)

				span.SetTag("requestID", reqID)

				defer func() {
					span.Finish()
				}()

				// nolint: errcheck
				tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

				r = r.WithContext(opentracing.ContextWithSpan(r.Context(), span))
			}

			next.ServeHTTP(w, r)
		})
	}
}
