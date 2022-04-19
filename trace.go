package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer wraps trace.Tracer to support goroutine-bound context propagation.
type Tracer interface {
	trace.Tracer

	// Continue is like Start, but starts the new span using the context from CurrentContext().
	Continue(spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

func NewTracer(name string, opts ...trace.TracerOption) Tracer {
	t := otel.Tracer(name, opts...)
	return &tracer{t}
}

type tracer struct {
	trace.Tracer
}

func (t *tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	newCtx, span := t.Tracer.Start(ctx, spanName, opts...)
	return newCtx, &spanWrapper{
		Span: span,
		ctx:  enterContext(newCtx),
	}
}

func (t *tracer) Continue(spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.Start(CurrentContext(), spanName, opts...)
}

type spanWrapper struct {
	trace.Span
	ctx *contextWrapper
}

func (s *spanWrapper) End(options ...trace.SpanEndOption) {
	s.Span.End(options...)
	s.ctx.exit()
}
