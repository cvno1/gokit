package log

import (
	"go.uber.org/zap/zapcore"
)

type options struct {
	level             zapcore.Level
	fields            []Field
	outputFile        string
	errOutputFile     string
	development       bool
	disableConsle     bool
	disableStackTrace bool
	disableCaller     bool
}

type Option func(*options)

// newOptions return a defaultoptions
func newOptions() *options {
	return &options{
		level:             zapcore.InfoLevel,
		fields:            []Field{},
		outputFile:        "stdout",
		errOutputFile:     "stderr",
		development:       false,
		disableStackTrace: false,
		disableCaller:     false,
	}
}

// A Level is a logging priority. Higher levels are more important.
type Level string

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = "debug"
	// InfoLevel is the default logging priority.
	InfoLevel = "info"
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = "warn"
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = "error"
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = "dpanic"
	// PanicLevel logs a message, then panics.
	PanicLevel = "panic"
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

// parseLevel parse string level to zapcore.Level
func parseLevel(level Level) zapcore.Level {
	zapcorLevel, _ := zapcore.ParseLevel(string(level))
	return zapcorLevel
}

// WithLevel set the default output level
func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = parseLevel(level)
	}
}

// WithField add some field(s) to log
func WithField(key string, value any) Option {
	return func(o *options) {
		o.fields = append(o.fields, NewField(key, value))
	}
}

// WithFile write log in some files
func WithFile(path string, errPath string) Option {
	return func(o *options) {
		o.outputFile = path
		o.errOutputFile = errPath
	}
}

// WithDevelopment set the env is development,
// which will changes the behavior of DPanicLevel and takes stacktraces more liberally
func WithDevelopment() Option {
	return func(o *options) {
		o.development = true
	}
}

// WithDisableCaller Disable output of caller information in the log
func WithDisableCaller() Option {
	return func(o *options) {
		o.disableCaller = true
	}
}

// WithDisableStackTrace disable the log to record a stack trace for
// all messages at or above panic level
func WithDisableStackTrace() Option {
	return func(o *options) {
		o.disableStackTrace = true
	}
}
