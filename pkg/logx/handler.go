package logx

import (
	"context"
	"log/slog"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

// censoringHandler wraps another slog.Handler to provide censoring of sensitive data.
// It processes log records by redacting sensitive information before passing them to
// the underlying handler.
type censoringHandler struct {
	// handler is the underlying slog.Handler that processes censored records
	handler slog.Handler
	// withDebug enables debug mode which bypasses censoring
	withDebug bool
	// sensitiveKeys defines the set of field names to be redacted
	sensitiveKeys map[string]struct{}
}

// Enabled reports whether the handler handles records at the given level.
// The implementation delegates to the underlying handler's Enabled method.
func (h *censoringHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler interface. It processes the log record by
// censoring sensitive attributes before passing to the underlying handler.
//
// The process maintains immutability of the original record by:
//   - Creating a new slice for censored attributes
//   - Processing each attribute through the censor
//   - Creating a new record with censored attributes
//   - Forwarding the censored record to the underlying handler
func (h *censoringHandler) Handle(ctx context.Context, r slog.Record) error {

	ctxAttrs := getAttrsFromContext(ctx)
	totalCap := r.NumAttrs()
	if ctxAttrs != nil {
		totalCap += len(ctxAttrs)
	}

	newAttrs := make([]slog.Attr, 0, totalCap)
	// Add context attributes first
	if ctxAttrs != nil {
		for _, attr := range ctxAttrs {
			if attr.Value.Kind() == slog.KindGroup {
				groupAttrs := attr.Value.Group()
				newGroupAttrs := make([]slog.Attr, 0, len(groupAttrs))

				for _, groupAttr := range groupAttrs {
					newGroupAttrs = append(newGroupAttrs, h.censorAttribute(groupAttr))
				}

				anyGroupAttrs := make([]any, len(newGroupAttrs))
				for i, groupAttr := range newGroupAttrs {
					anyGroupAttrs[i] = groupAttr
				}
				newAttrs = append(newAttrs, slog.Group(attr.Key, anyGroupAttrs...))
			} else {
				newAttrs = append(newAttrs, h.censorAttribute(attr))
			}
		}
	}

	r.Attrs(func(attr slog.Attr) bool {
		// Handle both direct attributes and groups
		if attr.Value.Kind() == slog.KindGroup {
			// Process group attributes
			groupAttrs := attr.Value.Group()
			newGroupAttrs := make([]slog.Attr, 0, len(groupAttrs))

			// Censor each attribute in the group
			for _, groupAttr := range groupAttrs {
				newGroupAttrs = append(newGroupAttrs, h.censorAttribute(groupAttr))
			}

			// Create new group with censored attributes
			anyGroupAttrs := make([]any, len(newGroupAttrs))
			for i, groupAttr := range newGroupAttrs {
				anyGroupAttrs[i] = groupAttr
			}
			newAttrs = append(newAttrs, slog.Group(attr.Key, anyGroupAttrs...))
		} else {
			newAttrs = append(newAttrs, h.censorAttribute(attr))
		}
		return true
	})

	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	for _, attr := range newAttrs {
		newRecord.AddAttrs(attr)
	}

	return h.handler.Handle(ctx, newRecord)
}

// WithAttrs implements slog.Handler interface. It returns a new Handler whose attributes
// consist of both the receiver's attributes and the arguments. The receiver's attributes
// appear first.
//
// The returned Handler maintains the same censoring configuration while wrapping a new
// underlying handler with the additional attributes.
func (h *censoringHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &censoringHandler{
		handler:       h.handler.WithAttrs(attrs),
		withDebug:     h.withDebug,
		sensitiveKeys: h.sensitiveKeys,
	}
}

// WithGroup implements slog.Handler interface. It returns a new Handler with the
// given group name added to its attributes. Groups are dot-separated names that
// appear in log output, depending on the output format.
//
// The returned Handler maintains the same censoring configuration while wrapping
// a new underlying handler with the specified group.
func (h *censoringHandler) WithGroup(name string) slog.Handler {
	return &censoringHandler{
		handler:       h.handler.WithGroup(name),
		withDebug:     h.withDebug,
		sensitiveKeys: h.sensitiveKeys,
	}
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

// GetAttrsFromContext retrieves the slog attributes stored in the context
func getAttrsFromContext(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		return attrs
	}
	return nil
}
