package logx

import (
	"context"
	"log/slog"
	"reflect"
)

// censorAttribute returns a new slog.Attr with sensitive data redacted.
// If debug mode is enabled (h.withDebug is true), it returns the original attribute unmodified.
// Otherwise, it performs deep cloning of the attribute value while redacting sensitive data
// based on the configured sensitiveKeys.
//
// The redaction process handles various data types including:
//   - Basic types (strings, numbers, etc.)
//   - Complex types (structs, maps, slices)
//   - JSON-encoded strings that may contain sensitive data
//
// For JSON strings, it attempts to parse and redact sensitive fields within the JSON
// before re-encoding.
func (h *censoringHandler) censorAttribute(attr slog.Attr) slog.Attr {
	if h.withDebug {
		return attr
	}
	ctx := context.Background()
	masked := clone(ctx, attr.Key, reflect.ValueOf(attr.Value.Any()), "", h.sensitiveKeys)
	return slog.Any(attr.Key, masked.Interface())
}
