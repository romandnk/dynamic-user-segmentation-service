package zap_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Config struct {
	Level           zapcore.Level
	Encoding        string
	OutputPath      []string
	ErrorOutputPath []string
}

type ZapLogger struct {
	Log *zap.Logger
}

func NewZapLogger(config Config) (*ZapLogger, error) {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(config.Level),
		Encoding:         config.Encoding,
		OutputPaths:      config.OutputPath,
		ErrorOutputPaths: config.ErrorOutputPath,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "lvl",
			TimeKey:    "ts",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.DateTime))
			},
			EncodeLevel: func(lvl zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(lvl.String())
			},
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{Log: logger}, nil
}

func (l *ZapLogger) Info(msg string, fields ...any) {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, field := range fields {
		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)
		default:
			return
		}
	}

	l.Log.Info(msg, zapFields...)
}

func (l *ZapLogger) Error(msg string, fields ...any) {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, field := range fields {
		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)
		default:
			return
		}
	}

	l.Log.Error(msg, zapFields...)
}
