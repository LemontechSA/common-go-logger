package logger

import (
	"context"

	"go.uber.org/zap"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type contextKey string

func (c contextKey) String() string {
	return "logger" + string(c)
}

// Customs context keys to avoid collisions
var (
	ContextKeyCorrelationID = contextKey("correlation_id")
	ContextKeyCausationID   = contextKey("causation_id")
	ContextKeyTenant        = contextKey("tenant")
	ContextKeyUserID        = contextKey("user_id")
	ContextKeyConsumer      = contextKey("consumer")
)

// Gets the correlation id from context
func getCorrelationID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(ContextKeyCorrelationID).(string)

	return value, ok
}

// Gets the causation id from context
func getCausationID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(ContextKeyCausationID).(string)

	return value, ok
}

// Gets the tenant id from context
func getTenant(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(ContextKeyTenant).(string)

	return value, ok
}

// Gets the user id from context
func getUserID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(ContextKeyUserID).(string)

	return value, ok
}

// Gets the consumer from context
func getConsumer(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(ContextKeyConsumer).(string)

	return value, ok
}

// Gets the datadog trace id from context
func getDDTraceID(ctx context.Context) (string, bool) {
	span, _ := ddtracer.SpanFromContext(ctx)
	value := convertUInt64ToString(span.Context().TraceID())

	if value != "" && value != "0" {
		return value, true
	}

	return value, false
}

// Gets the datadog span id from context
func getDDSpanID(ctx context.Context) (string, bool) {
	span, _ := ddtracer.SpanFromContext(ctx)
	value := convertUInt64ToString(span.Context().SpanID())

	if value != "" && value != "0" {
		return value, true
	}

	return value, false
}

// Puts the values of the context into the log
func withContext(ctx context.Context, l *contextLogger) *zap.Logger {
	logger := l.logger.WithOptions()
	ddValues := map[string]string{}

	if correlationID, exists := getCorrelationID(ctx); exists {
		logger = logger.With(zap.String("correlation_id", correlationID))
	}

	if causationID, exists := getCausationID(ctx); exists {
		logger = logger.With(zap.String("causation_id", causationID))
	}

	if tenant, exists := getTenant(ctx); exists {
		logger = logger.With(zap.String("tenant", tenant))
	}

	if userID, exists := getUserID(ctx); exists {
		logger = logger.With(zap.String("user_id", userID))
	}

	if consumer, exists := getConsumer(ctx); exists {
		logger = logger.With(zap.String("consumer", consumer))
	}

	if ddTraceID, exists := getDDTraceID(ctx); exists {
		ddValues["trace_id"] = ddTraceID
	}

	if ddSpanID, exists := getDDSpanID(ctx); exists {
		ddValues["span_id"] = ddSpanID
	}

	if len(ddValues) > 0 {
		if l.config.Version != "" {
			ddValues["version"] = l.config.Version
		}

		if l.config.Environment != "" {
			ddValues["env"] = l.config.Environment
		}

		if l.config.Service != "" {
			ddValues["service"] = l.config.Service
		}

		logger = logger.With(zap.Any("dd", ddValues), zap.String("ddsource", "go"))
	}

	return logger
}
