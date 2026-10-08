package logger

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"testing"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// captureStdout swaps os.Stdout for a pipe while newLogger runs, so the
// logger it builds writes to the pipe. It returns the decoded JSON lines.
func captureStdout(t *testing.T, newLogger func(), log func()) []map[string]interface{} {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	original := os.Stdout
	os.Stdout = writer
	newLogger()
	os.Stdout = original

	log()
	writer.Close()

	var lines []map[string]interface{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var line map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			t.Fatalf("log line is not JSON: %q: %v", scanner.Text(), err)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		t.Fatalf("read stdout: %v", err)
	}

	return lines
}

func TestLoggerWritesExpectedFields(t *testing.T) {
	config := Configuration{
		Environment:  "test",
		Service:      "svc",
		Team:         "team",
		Project:      "project",
		ConsoleLevel: "info",
		Version:      "1.2.3",
	}

	var l Logger
	lines := captureStdout(t,
		func() { l = NewLogger(config) },
		func() {
			l.Debug("dropped", "", nil)
			l.Info("hello", "do-thing", map[string]string{"duration": "42", "key": "value"})
		},
	)

	if len(lines) != 1 {
		t.Fatalf("got %d log lines, want 1 (debug must be filtered at info level): %v", len(lines), lines)
	}

	line := lines[0]
	for key, want := range map[string]string{
		"level":       "info",
		"message":     "hello",
		"action":      "do-thing",
		"service":     "svc",
		"environment": "test",
		"team":        "team",
		"project":     "project",
		"version":     "1.2.3",
	} {
		if got := line[key]; got != want {
			t.Errorf("field %q = %v, want %q", key, got, want)
		}
	}

	for _, key := range []string{"date", "method", "pid", "host"} {
		if _, ok := line[key]; !ok {
			t.Errorf("field %q missing", key)
		}
	}

	payload, ok := line["payload"].(map[string]interface{})
	if !ok {
		t.Fatalf("payload = %v, want an object", line["payload"])
	}
	if payload["duration"] != float64(42) {
		t.Errorf("payload.duration = %v, want number 42", payload["duration"])
	}
	if payload["key"] != "value" {
		t.Errorf("payload.key = %v, want %q", payload["key"], "value")
	}
}

func TestContextLoggerAddsContextAndTraceFields(t *testing.T) {
	tracer := mocktracer.Start()
	defer tracer.Stop()

	span, ctx := ddtracer.StartSpanFromContext(context.Background(), "op")
	defer span.Finish()

	ctx = context.WithValue(ctx, ContextKeyCorrelationID, "corr")
	ctx = context.WithValue(ctx, ContextKeyCausationID, "cause")
	ctx = context.WithValue(ctx, ContextKeyTenant, "tenant")
	ctx = context.WithValue(ctx, ContextKeyUserID, "user")
	ctx = context.WithValue(ctx, ContextKeyConsumer, "consumer")

	config := Configuration{Environment: "test", Service: "svc", Version: "1.2.3"}

	var l ContextLogger
	lines := captureStdout(t,
		func() { l = NewContextLogger(config) },
		func() { l.CtxInfo(ctx, "hello", "", nil) },
	)

	if len(lines) != 1 {
		t.Fatalf("got %d log lines, want 1: %v", len(lines), lines)
	}

	line := lines[0]
	for key, want := range map[string]string{
		"correlation_id": "corr",
		"causation_id":   "cause",
		"tenant":         "tenant",
		"user_id":        "user",
		"consumer":       "consumer",
		"ddsource":       "go",
	} {
		if got := line[key]; got != want {
			t.Errorf("field %q = %v, want %q", key, got, want)
		}
	}

	dd, ok := line["dd"].(map[string]interface{})
	if !ok {
		t.Fatalf("dd = %v, want an object", line["dd"])
	}
	for key, want := range map[string]string{
		"trace_id": strconv.FormatUint(span.Context().TraceID(), 10),
		"span_id":  strconv.FormatUint(span.Context().SpanID(), 10),
		"env":      "test",
		"service":  "svc",
		"version":  "1.2.3",
	} {
		if got := dd[key]; got != want {
			t.Errorf("dd.%s = %v, want %q", key, got, want)
		}
	}
}

func TestContextLoggerWithoutSpanOmitsTraceFields(t *testing.T) {
	var l ContextLogger
	lines := captureStdout(t,
		func() { l = NewContextLogger(Configuration{}) },
		func() { l.CtxWarn(context.Background(), "hello", "", nil) },
	)

	if len(lines) != 1 {
		t.Fatalf("got %d log lines, want 1: %v", len(lines), lines)
	}
	if _, ok := lines[0]["dd"]; ok {
		t.Errorf("dd field present without a span: %v", lines[0]["dd"])
	}
	if _, ok := lines[0]["ddsource"]; ok {
		t.Errorf("ddsource field present without a span")
	}
}
