package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the contract for the logger
type Logger interface {
	Debug(description string, message string, payload map[string]string)
	Info(description string, message string, payload map[string]string)
	Warn(description string, message string, payload map[string]string)
	Error(description string, message string, payload map[string]string)
	Fatal(description string, message string, payload map[string]string)
	SetContext(ctx context.Context)
	ClearContext()
}

// Configuration stores the config for the logger
type Configuration struct {
	Environment  string
	Service      string
	Team         string
	Project      string
	ConsoleLevel string
}

type zapLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

// NewLogger returns an instance of logger
func NewLogger(config Configuration) Logger {
	ctx := context.TODO()
	level := getLogLevel(config.ConsoleLevel)
	writer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(getEncoder(), writer, level)
	log := zap.New(
		core,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	)
	logger := createLogger(config, log)

	return &zapLogger{logger, ctx}
}

func (l *zapLogger) Debug(description string, message string, payload map[string]string) {
	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		l.logger.Debug(
			message,
			zap.String("description", description),
			zap.Any("payload", parsedPayload),
		)
	} else {
		l.logger.Debug(
			message,
			zap.String("description", description),
		)
	}
}

func (l *zapLogger) Info(description string, message string, payload map[string]string) {
	log := withContext(l.ctx, l.logger)

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		log.Info(
			message,
			zap.String("description", description),
			zap.Any("payload", parsedPayload),
		)
	} else {
		log.Info(
			message,
			zap.String("description", description),
		)
	}
}

func (l *zapLogger) Warn(description string, message string, payload map[string]string) {
	log := withContext(l.ctx, l.logger)

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		log.Warn(
			message,
			zap.String("description", description),
			zap.Any("payload", parsedPayload),
		)
	} else {
		log.Warn(
			message,
			zap.String("description", description),
		)
	}
}

func (l *zapLogger) Error(description string, message string, payload map[string]string) {
	log := withContext(l.ctx, l.logger)

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		log.Error(
			message,
			zap.String("description", description),
			zap.Any("payload", parsedPayload),
		)
	} else {
		log.Error(
			message,
			zap.String("description", description),
		)
	}
}

func (l *zapLogger) Fatal(description string, message string, payload map[string]string) {
	log := withContext(l.ctx, l.logger)

	if len(payload) > 0 {
		parsedPayload := parserPayload(payload)

		log.Fatal(
			message,
			zap.String("description", description),
			zap.Any("payload", parsedPayload),
		)
	} else {
		log.Fatal(
			message,
			zap.String("description", description),
		)
	}
}

func (l *zapLogger) SetContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *zapLogger) ClearContext() {
	l.ctx = context.TODO()
}
