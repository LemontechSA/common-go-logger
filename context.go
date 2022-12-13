package logger

import (
	"context"
	"net/http"
)

var (
	contextKeyTraceID      = contextKey("trace_id")
	contextKeyRequestID    = contextKey("request_id")
	contextKeySessionID    = contextKey("session_id")
	contextKeyConsumerName = contextKey("consumer_name")
)

type contextKey string

func (c contextKey) String() string {
	return "logger" + string(c)
}

// traceID gets the trace id from context
func traceID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKeyTraceID).(string)

	return id, ok
}

// requestID gets the request id from context
func requestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKeyRequestID).(string)

	return id, ok
}

// sessionID gets the session id from context
func sessionID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKeySessionID).(string)

	return id, ok
}

// consumerName gets the consumer name from context
func consumerName(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKeyConsumerName).(string)

	return id, ok
}

// CreateRequestContext creates the context from request
func CreateRequestContext(req *http.Request) context.Context {
	ctx := req.Context()

	var (
		traceID      = req.Header.Get("trace_id")
		requestID    = req.Header.Get("request_id")
		sessionID    = req.Header.Get("session_id")
		consumerName = req.Header.Get("consumer_name")
	)

	if traceID != "" {
		ctx = context.WithValue(ctx, contextKeyTraceID, traceID)
	}

	if requestID != "" {
		ctx = context.WithValue(ctx, contextKeyRequestID, requestID)
	}

	if sessionID != "" {
		ctx = context.WithValue(ctx, contextKeySessionID, sessionID)
	}

	if consumerName != "" {
		ctx = context.WithValue(ctx, contextKeyConsumerName, consumerName)
	}

	return ctx
}
