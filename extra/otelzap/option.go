package otelzap

import "go.uber.org/zap/zapcore"

// Option applies a configuration to the given config.
type Option func(l *Logger)

// WithMinLevel sets the minimal zap logging level on which the log message
// is recorded on the span.
//
// The default is >= zap.WarnLevel.
func WithMinLevel(lvl zapcore.Level) Option {
	return func(l *Logger) {
		l.minLevel = lvl
	}
}

// WithErrorStatusLevel sets the minimal zap logging level on which
// the span status is set to codes.Error.
//
// The default is >= zap.ErrorLevel.
func WithErrorStatusLevel(lvl zapcore.Level) Option {
	return func(l *Logger) {
		l.errorStatusLevel = lvl
	}
}

// WithCaller configures the logger to annotate each event with the filename,
// line number, and function name of the caller.
//
// It is enabled by default.
func WithCaller(flag bool) Option {
	return func(l *Logger) {
		l.caller = flag
	}
}

// WithStackTrace configures the logger to capture logs with a stack trace.
func WithStackTrace(flag bool) Option {
	return func(l *Logger) {
		l.stackTrace = flag
	}
}
