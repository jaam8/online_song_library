package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

//const Key = "logger"

//type Logger struct {
//	l *zap.Logger
//}

func New(logLevel string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "time"
	//config.EncoderConfig.CallerKey = ""

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05")) // Форматируем время как нужно
	})

	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	}
	config.Level.SetLevel(level)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, err
}

//func GetLoggerFromCtx(ctx context.Context) *Logger {
//	return ctx.Value(Key)
//}

//func (l *Logger) Info(msg string, fields ...zap.Field) {
//	l.l.Info(msg, fields...)
//}
//
//func (l *Logger) Fatal(msg string, fields ...zap.Field) {
//	l.l.Fatal(msg, fields...)
//}
//
//func (l *Logger) Debug(msg string, fields ...zap.Field) {
//	l.l.Debug(msg, fields...)
//}
