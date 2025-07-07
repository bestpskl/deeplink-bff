package logx

import (
	snake "deeplink-bff/pkg/string"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mdobak/go-xerrors"
)

// Config contains essential logger configuration like environment and application source.
type Config struct {
	Environment string `json:"env"`
	Source      string `json:"source"`
}

// options holds internal logger settings including redaction configuration.
type optionsConfig struct {
	Level                slog.Level
	AddSource            bool     `json:"add_source"`
	SensitiveKeys        []string `json:"sensitive_keys"`
	WithDebug            bool     `json:"with_debug"`
	DefaultRedactMessage string   `json:"default_redact_message"`
}

type stackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
}

const (
	// DefaultTagKey is a default key name of struct tag for sensitive.
	DefaultTagKey = "sensitive"
	// DefaultRedactMessage is a default message to replace redacted value.
	DefaultRedactMessage = "*"
)

var sensitivekeys = map[string]struct{}{}

// New creates a slog.Logger with redaction and structured logging capabilities.
//
// The logger outputs JSON-formatted logs with:
//   - Automatic redaction of sensitive fields
//   - Standard fields for environment and application name
//   - Runtime information including PID and Go version
//
// Example:
//
//	log, err := logx.New(logger.Config{
//	    Environment: "development",
//	    Source:     "myapp",
//	})
func New(cfg Config, options ...Option) (*slog.Logger, error) {

	// buildInfo := newBuildInfo()
	logOpts := optionsConfig{
		Level:                slog.LevelDebug,
		AddSource:            false,
		SensitiveKeys:        []string{},
		WithDebug:            false,
		DefaultRedactMessage: DefaultRedactMessage,
	}

	// Validate blacklist during package initialization
	// if err := validateBlackList(blackList); err != nil {
	// 	panic(fmt.Sprintf("invalid blacklist configuration: %v", err))
	// }

	for _, option := range options {
		if err := option(&logOpts); err != nil {
			return nil, fmt.Errorf("failed to apply logger option: %w", err)
		}
	}

	// Set sensitive keys
	sensitivekeys = make(map[string]struct{})
	for _, key := range logOpts.SensitiveKeys {
		sensitivekeys[snake.SnakeCase(key)] = struct{}{}
	}
	// Set sensitive keys from blacklist by default
	for key := range blackList {
		sensitivekeys[snake.SnakeCase(key)] = struct{}{}
	}

	opts := &slog.HandlerOptions{
		Level:       logOpts.Level,
		AddSource:   logOpts.AddSource,
		ReplaceAttr: replaceErrorAttribute,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Wrap with censoring handler
	censoringHandler := &censoringHandler{
		handler:       handler,
		withDebug:     logOpts.WithDebug,
		sensitiveKeys: sensitivekeys,
	}

	// Create the logger with default fields
	logger := slog.New(censoringHandler).With(
		slog.String("source", cfg.Source),
		slog.String("env", cfg.Environment),
		// slog.Group("program_info",
		// 	slog.Int("pid", os.Getpid()),
		// 	slog.String("go_version", buildInfo.GoVersion),
		// ),
	)

	return logger, nil
}
func replaceErrorAttribute(groups []string, attr slog.Attr) slog.Attr {
	switch attr.Value.Kind() {
	case slog.KindAny:
		switch v := attr.Value.Any().(type) {
		case error:
			attr.Value = formatErrorWithTrace(v)
		}
	}

	return attr
}

// extractStackFrames extracts and formats stack frames from an error
func extractStackFrames(err error) []stackFrame {
	trace := xerrors.StackTrace(err)

	if len(trace) == 0 {
		return nil
	}

	frames := trace.Frames()
	formattedFrames := make([]stackFrame, len(frames))

	for i, frame := range frames {
		formattedFrames[i] = stackFrame{
			Source: filepath.Join(
				filepath.Base(filepath.Dir(frame.File)),
				filepath.Base(frame.File),
			),
			Func: filepath.Base(frame.Function),
			Line: frame.Line,
		}
	}

	return formattedFrames
}

// formatErrorWithTrace returns a slog.Value containing the error message and stack trace.
// If the error doesn't implement the StackTrace interface, only the message is included.
func formatErrorWithTrace(err error) slog.Value {
	var attributes []slog.Attr

	attributes = append(attributes, slog.String("msg", err.Error()))

	frames := extractStackFrames(err)
	if frames != nil {
		attributes = append(attributes,
			slog.Any("trace", frames),
		)
	}

	return slog.GroupValue(attributes...)
}
