package log

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Dpanic(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	Flush()
}

type zapLogger struct {
	logger *zap.Logger
}

var _ Logger = (*zapLogger)(nil)

func New(options ...Option) Logger {
	opts := newOptions()

	for _, o := range options {
		o(opts)
	}

	return &zapLogger{newZapLogger(opts)}
}

func newZapLogger(o *options) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendFloat64(float64(d) / float64(time.Millisecond))
		},
		EncodeCaller:        zapcore.FullCallerEncoder,
		EncodeName:          zapcore.FullNameEncoder,
		NewReflectedEncoder: nil, // default json encoder
		ConsoleSeparator:    "\t",
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if o.development {
			return lvl >= zapcore.ErrorLevel
		} else {
			return lvl >= o.level && lvl >= zapcore.ErrorLevel
		}
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if o.development {
			return lvl < zapcore.ErrorLevel
		} else {
			return lvl >= o.level && lvl < zapcore.ErrorLevel
		}
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	consoleCores := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), stdout, lowPriority),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), stderr, highPriority),
	)

	preparedir := func(file string) {
		dir := filepath.Dir(file)
		if err := os.MkdirAll(dir, 0766); err != nil {
			panic(err)
		}
	}

	preparedir(o.outputFile)
	preparedir(o.errOutputFile)

	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	fileCores := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   o.outputFile,
				MaxSize:    128,
				MaxAge:     30,
				MaxBackups: 300,
				LocalTime:  true,
				Compress:   true,
			}),
			lowPriority,
		),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   o.errOutputFile,
				MaxSize:    128,
				MaxAge:     30,
				MaxBackups: 300,
				LocalTime:  true,
				Compress:   true,
			}),
			highPriority,
		),
	)
	var zopts []zap.Option

	if o.development {
		zopts = append(zopts, zap.Development())
	}
	if !o.disableCaller {
		zopts = append(zopts, zap.AddCaller())
	}
	if !o.disableStackTrace {
		zopts = append(zopts, zap.AddStacktrace(highPriority))
	}

	zopts = append(zopts, zap.Fields(toZapField(o.fields...)...))

	logger := zap.New(zapcore.NewTee(consoleCores, fileCores),
		zopts...,
	)
	return logger
}

func (z *zapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, toZapField(fields...)...)
}

func (z *zapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, toZapField(fields...)...)
}

func (z *zapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, toZapField(fields...)...)
}

func (z *zapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, toZapField(fields...)...)
}

func (z *zapLogger) Dpanic(msg string, fields ...Field) {
	z.logger.DPanic(msg, toZapField(fields...)...)
}

func (z *zapLogger) Panic(msg string, fields ...Field) {
	z.logger.Panic(msg, toZapField(fields...)...)
}

func (z *zapLogger) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, toZapField(fields...)...)
}

func (z *zapLogger) Flush() {
	_ = z.logger.Sync()
}
