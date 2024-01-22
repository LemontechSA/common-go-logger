package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the contract for the logger
type Logger interface {
	Debug(message string, action string, payload map[string]string)
	Info(message string, action string, payload map[string]string)
	Warn(message string, action string, payload map[string]string)
	Error(message string, action string, payload map[string]string)
	Fatal(message string, action string, payload map[string]string)
}

// ContextLogger is the contract for the context logger
type ContextLogger interface {
	CtxDebug(ctx context.Context, message string, action string, payload map[string]string)
	CtxInfo(ctx context.Context, message string, action string, payload map[string]string)
	CtxWarn(ctx context.Context, message string, action string, payload map[string]string)
	CtxError(ctx context.Context, message string, action string, payload map[string]string)
	CtxFatal(ctx context.Context, message string, action string, payload map[string]string)
}

// Configuration stores the config for the logger
type Configuration struct {
	Environment  string
	Service      string
	Team         string
	Project      string
	ConsoleLevel string
	Version      string
}

type genericLogger struct {
	logger *zap.Logger
}

type contextLogger struct {
	logger *zap.Logger
	config Configuration
}

// NewLogger returns an instance of logger
func NewLogger(config Configuration) Logger {
	logger := createLogger(config)

	return &genericLogger{logger}
}

// NewContextLogger returns an instance of context logger
func NewContextLogger(config Configuration) ContextLogger {
	logger := createLogger(config)

	return &contextLogger{logger, config}
}

func (l *genericLogger) Debug(message string, action string, payload map[string]string) {
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	l.logger.Debug(message, fields...)
}

func (l *genericLogger) Info(message string, action string, payload map[string]string) {
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	l.logger.Info(message, fields...)
}

func (l *genericLogger) Warn(message string, action string, payload map[string]string) {
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	l.logger.Warn(message, fields...)
}

func (l *genericLogger) Error(message string, action string, payload map[string]string) {
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	l.logger.Error(message, fields...)
}

func (l *genericLogger) Fatal(message string, action string, payload map[string]string) {
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	l.logger.Fatal(message, fields...)
}

func (l *contextLogger) CtxDebug(ctx context.Context, message string, action string, payload map[string]string) {
	log := withContext(ctx, l)
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	log.Debug(message, fields...)
}

func (l *contextLogger) CtxInfo(ctx context.Context, message string, action string, payload map[string]string) {
	log := withContext(ctx, l)
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	log.Info(message, fields...)
}

func (l *contextLogger) CtxWarn(ctx context.Context, message string, action string, payload map[string]string) {
	log := withContext(ctx, l)
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	log.Warn(message, fields...)
}

func (l *contextLogger) CtxError(ctx context.Context, message string, action string, payload map[string]string) {
	log := withContext(ctx, l)
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	log.Error(message, fields...)
}

func (l *contextLogger) CtxFatal(ctx context.Context, message string, action string, payload map[string]string) {
	log := withContext(ctx, l)
	fields := []zapcore.Field{}

	if action != "" {
		fields = append(fields, zap.String("action", action))
	}

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		fields = append(fields, zap.Any("payload", parsedPayload))
	}

	log.Fatal(message, fields...)
}
