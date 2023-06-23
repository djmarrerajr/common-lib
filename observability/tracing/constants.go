package tracing

// nolint: unused
const (
	DefaultTracingSamplerType  = "probabilistic"
	DefaultTracingSamplerValue = float64(.50)
)

// nolint: unused
const (
	TracingHostPortEnvKey     = "TRACING_HOST_AND_PORT"
	TracingDisabledEnvKey     = "TRACING_DISABLED"
	TracingLogSpansEnvKey     = "TRACING_LOG_SPANS"
	TracingSamplerTypeEnvKey  = "TRACING_SAMPLER_TYPE"
	TracingSamplerValueEnvKey = "TRACING_SAMPLER_VALUE"
)

// nolint: unused
const (
	HeaderRequestId = "X-Request-Id"
)
