package logx

import (
	"io"
	"log/slog"
)

// Option is a functional option for configuring the logger.
type Option func(*optionsConfig) error

// WithLevel sets the log level for the logger.
// The log level determines the severity of messages that will be logged.
func WithLevel(level slog.Level) Option {
	return func(cfg *optionsConfig) error {
		cfg.Level = level
		return nil
	}
}

// WithAddSource enables or disables the inclusion of source information in logs.
// When enabled, log entries may include details such as filename and line number.
func WithAddSource(addSource bool) Option {
	return func(cfg *optionsConfig) error {
		cfg.AddSource = addSource
		return nil
	}
}

// WithSensitiveKeys specifies keys that should be treated as sensitive information.
// These keys may be redacted or masked in logs to prevent leakage of sensitive data.
func WithSensitiveKeys(keys []string) Option {
	return func(cfg *optionsConfig) error {
		cfg.SensitiveKeys = keys
		return nil
	}
}

// WithDebug enables or disables debug mode for the logger.
// When enabled, additional debugging information may be included in log output.
func WithDebugMode(withDebug bool) Option {
	return func(cfg *optionsConfig) error {
		cfg.WithDebug = withDebug
		return nil
	}
}

// WithDefaultRedactMessage sets the default message to use when redacting sensitive data.
func WithDefaultRedactMessage(defaultRedactMessage string) Option {
	return func(cfg *optionsConfig) error {
		cfg.DefaultRedactMessage = defaultRedactMessage
		return nil
	}
}

// WithWriter sets a custom io.Writer for logger output (e.g., file or MultiWriter).
func WithWriter(w io.Writer) Option {
	return func(cfg *optionsConfig) error {
		cfg.Writer = w
		return nil
	}
}
