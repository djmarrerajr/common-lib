package traces

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/djmarrerajr/common-lib/errs"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

func NewTracerFromEnv(env utils.Environ, appCtx shared.ApplicationContext, appName, appVersion string) (opentracing.Tracer, io.Closer, error) {
	agent, err := env.GetRequired(TracingHostPortEnvKey)
	if err != nil {
		return nil, nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	samplerType, OK := env.Get(TracingSamplerTypeEnvKey)
	if !OK {
		samplerType = DefaultTracingSamplerType
	}

	enabled, _, err := env.GetBool(TracingDisabledEnvKey)
	if err != nil {
		return nil, nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	logSpans, _, err := env.GetBool(TracingLogSpansEnvKey)
	if err != nil {
		return nil, nil, errs.WithType(err, errs.ErrTypeConfiguration)
	}

	cfg := &config.Configuration{
		ServiceName: appName,
		Disabled:    enabled,
		Sampler: &config.SamplerConfig{
			Type:  samplerType,
			Param: DefaultTracingSamplerValue,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           logSpans,
			LocalAgentHostPort: agent,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err == nil {
		opentracing.SetGlobalTracer(tracer)
	}

	return tracer, closer, err
}
