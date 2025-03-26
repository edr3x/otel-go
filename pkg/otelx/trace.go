package otelx

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// pass the key value pair when needed like
// example:
// oteltrace.WithAttributes(attribute.String("id", id))
func StartSpan(ctx context.Context, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	pc, _, line, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	return tracer.Start(ctx, fmt.Sprintf("%s:%d", details.Name(), line), opts...)
}
