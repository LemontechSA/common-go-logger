package logger

import (
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
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func createLogger(config Configuration) *zap.Logger {
	level := getLogLevel(config.ConsoleLevel)
	writer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(getEncoder(), writer, level)
	logger := zap.New(
		core,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	)
	host, _ := os.Hostname()
	pid := os.Getpid()

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

	if config.Version != "" {
		logger = logger.With(zap.String("version", config.Version))
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

func convertUInt64ToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
