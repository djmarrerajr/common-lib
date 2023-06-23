package tracing

import (
	"context"
	"runtime"

	"github.com/opentracing/opentracing-go"
)

func StartChildSpan(ctx context.Context, name string) (opentracing.Span, context.Context) {
	childSpan, _ := opentracing.StartSpanFromContext(ctx, name)

	// grab some information re: our caller and add it to the span
	currFrame := getCurrStackFrame()

	childSpan.SetTag("file", currFrame.File)
	childSpan.SetTag("line", currFrame.Line)
	childSpan.SetTag("func", currFrame.Function)

	return childSpan, opentracing.ContextWithSpan(ctx, childSpan)
}

func FinishChildSpan(childSpan opentracing.Span) {
	childSpan.Finish()
}

func getCurrStackFrame() runtime.Frame {
	callers := make([]uintptr, 10)
	runtime.Callers(3, callers)

	callFrames := runtime.CallersFrames(callers)
	currFrame, _ := callFrames.Next()

	return currFrame
}
