package traces

import (
	"io"

	"github.com/opentracing/opentracing-go"

	"github.com/djmarrerajr/common-lib/utils"
)

func NewTracerFromEnv(env utils.Environ, appName, appVersion string) (opentracing.Tracer, io.Closer, error) {
	// agent, err := env.GetRequired(TracingHostPortEnvKey)
	// if err != nil {
	// 	return nil, nil, errs.WithType(err, errs.ErrTypeConfiguration)
	// }

	// samplerType, OK := env.Get(TracingSamplerTypeEnvKey)
	// if !OK {
	// 	samplerType = DefaultTracingSamplerType
	// }

	// enabled, _, err := env.GetBool(TracingDisabledEnvKey)
	// if err != nil {
	// 	return nil, nil, errs.WithType(err, errs.ErrTypeConfiguration)
	// }

	// logSpans, _, err := env.GetBool(TracingLogSpansEnvKey)
	// if err != nil {
	// 	return nil, nil, errs.WithType(err, errs.ErrTypeConfiguration)
	// }

	// tracerConfig := &tracing.TracerConfig{
	// 	AgentHostPort: agent,
	// 	Disabled:      enabled,
	// 	Environment:   "dev",
	// 	LogSpans:      logSpans,
	// 	SamplerType:   tracing.SamplerType(samplerType),
	// 	SamplerValue:  DefaultTracingSamplerValue,
	// 	ServiceName:   appName,
	// 	Version:       appVersion,
	// }

	return nil, nil, nil
}
