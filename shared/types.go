package shared

import (
	"context"
	"io"

	"github.com/djmarrerajr/common-lib/observability/metrics"
	"github.com/djmarrerajr/common-lib/services/db"
	"github.com/djmarrerajr/common-lib/utils"
	"github.com/go-playground/validator"
	"github.com/opentracing/opentracing-go"
)

type ApplicationContext struct {
	RootCtx   context.Context     // root context
	Logger    utils.Logger        // root logger for the application
	Tracer    opentracing.Tracer  // tracing service (i.e. Jaeger)
	Collector metrics.Collector   // metrics collector (i.e. Prometheus)
	Validator *validator.Validate // struct validator
	Server    Servable            // embedded HTTP/HTTPS server
	Database  db.Adapter          // embedded Database adapter

	Closer io.Closer
}
