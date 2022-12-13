package logger

import (
	"context"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05.000Z"))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "method"
	encoderConfig.TimeKey = "date"
	encoderConfig.EncodeTime = timeEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func createLogger(config Configuration, logger *zap.Logger) *zap.Logger {
	host, _ := os.Hostname()
	pid := os.Getegid()

	logger = logger.With(
		zap.Int("pid", pid),
		zap.String("host", host),
	)

	if config.Service != "" {
		logger = logger.With(zap.String("service", config.Service))
	}

	if config.Environment != "" {
		logger = logger.With(zap.String("environment", config.Environment))
	}

	if config.Team != "" {
		logger = logger.With(zap.String("team", config.Team))
	}

	if config.Project != "" {
		logger = logger.With(zap.String("project", config.Project))
	}

	return logger
}

func withContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if tID, exists := traceID(ctx); exists {
		logger = logger.With(zap.String("trace_id", tID))
	}

	if rID, exists := requestID(ctx); exists {
		logger = logger.With(zap.String("request_id", rID))
	}

	if sID, exists := sessionID(ctx); exists {
		logger = logger.With(zap.String("session_id", sID))
	}

	if cName, exists := consumerName(ctx); exists {
		logger = logger.With(zap.String("consumer_name", cName))
	}

	return logger
}

func parserPayload(payload map[string]string) map[string]interface{} {
	parsedPayload := make(map[string]interface{})

	for key, value := range payload {
		if key == "duration" {
			newValue, _ := strconv.Atoi(value)

			parsedPayload[key] = newValue
		} else {
			parsedPayload[key] = value
		}
	}

	return parsedPayload
}
