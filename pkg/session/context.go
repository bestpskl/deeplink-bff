package session

import (
	"context"
)

// sessionInfoContextKey is the type used as a context key to prevent collisions
type sessionInfoContextKey struct{}

// defaultSessionInfoContextKey is the singleton key instance for session storage
var defaultSessionInfoContextKey = sessionInfoContextKey{}

// Get retrieves session information from the context.
// Returns the session Info and true if found, or nil and false if not present.
//
// Example:
//
//	info, ok := session.Get(ctx)
//	if !ok {
//	    return errors.New("no session found")
//	}
func Get(ctx context.Context) (*Info, bool) {
	// retrieve from default context key first if
	// session already been decoded
	if info, ok := ctx.Value(defaultSessionInfoContextKey).(*Info); ok {
		return info, true
	}

	return nil, false
}

// MustGet retrieves session information from context, panicking if not found.
// Use this when session presence is required for operation to continue.
//
// Example:
//
//	info := session.MustGet(ctx) // Panics if no session
func MustGet(ctx context.Context) *Info {
	if session, ok := Get(ctx); !ok {
		panic("unauthenticated")
	} else {
		return session
	}
}

// WithInfo creates a new context containing the provided session information.
// This is the primary way to store session data in a context.
//
// Example:
//
//	ctx = session.WithInfo(ctx, &session.Info{})
func WithInfo(ctx context.Context, info *Info) context.Context {
	return context.WithValue(ctx, defaultSessionInfoContextKey, info)
}

// AppendLanguage adds or updates the language setting in the session context.
// Creates a new session if one doesn't exist.
//
// Example:
//
//	ctx = session.AppendLanguage(ctx, "en-US")
func AppendLanguage(ctx context.Context, language string) context.Context {
	if info, ok := Get(ctx); ok {
		info.session.language = language
		return ctx
	}

	var info = &Info{}
	info.SetLanguage(language)
	return WithInfo(ctx, info)
}
