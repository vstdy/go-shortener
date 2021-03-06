package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/vstdy/go-shortener/pkg/logging"
)

type SpanOption func(span opentracing.Span)

// CustomComponentTagOpt sets component Span tag.
func CustomComponentTagOpt(value string) SpanOption {
	return func(span opentracing.Span) {
		ext.Component.Set(span, value)
	}
}

// CustomTagOpt sets custom Span tag.
func CustomTagOpt(key, value string) SpanOption {
	return func(span opentracing.Span) {
		span.SetTag(key, value)
	}
}

// CorrelationIDOpt adds correlation ID Span tag.
func CorrelationIDOpt(ctx context.Context) SpanOption {
	return func(span opentracing.Span) {
		correlationID, err := logging.GetCorrelationID(ctx)
		if err != nil {
			correlationID = "undefined"
		}
		CustomTagOpt(logging.CorrelationIDKey, correlationID)(span)
	}
}
