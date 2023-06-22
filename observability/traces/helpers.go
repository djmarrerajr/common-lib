package traces

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func StartChildSpan(ctx context.Context, name string) (opentracing.Span, context.Context) {
	childSpan, _ := opentracing.StartSpanFromContext(ctx, name)

	return childSpan, opentracing.ContextWithSpan(ctx, childSpan)
}

func FinishChildSpan(childSpan opentracing.Span) {
	childSpan.Finish()
}
