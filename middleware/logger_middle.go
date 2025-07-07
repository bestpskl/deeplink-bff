package middleware

import (
	"deeplink-bff/pkg/logx"
	"log/slog"
	"path/filepath"
	"runtime"

	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const (
	customAttributesCtxKey = "slog-fiber.custom-attributes"
	requestIDCtx           = "slog-fiber.request-id"
)

var (
	RequestIDKey = "request_id"
	TraceIDKey   = "trace_id"
	SpanIDKey    = "span_id"

	RequestBodyMaxSize  = 64 * 1024 // 64KB
	ResponseBodyMaxSize = 64 * 1024 // 64KB

	HiddenRequestHeaders = map[string]struct{}{
		"authorization": {},
		"cookie":        {},
		"set-cookie":    {},
		"x-auth-token":  {},
		"x-csrf-token":  {},
		"x-xsrf-token":  {},
	}
	HiddenResponseHeaders = map[string]struct{}{
		"set-cookie": {},
	}

	// Formatted with http.CanonicalHeaderKey
	RequestIDHeaderKey = "X-Request-Id"
)

// Config defines logging behavior for the middleware
type Config struct {
	DefaultLevel     slog.Level // Level for successful requests
	ClientErrorLevel slog.Level // Level for 4xx responses
	ServerErrorLevel slog.Level // Level for 5xx responses

	// Base Group
	WithUserAgent bool
	// Trace Group
	WithTraceID   bool
	WithSpanID    bool
	WithRequestID bool
	// Request Group
	WithRequestBody   bool
	WithRequestHeader bool
	// Response Group
	WithResponseBody   bool
	WithResponseHeader bool
	// AWS Group
	WithXRay bool
}

// Logger returns a Fiber middleware with default configuration.
func Logger() fiber.Handler {
	return LoggerWithConfig(Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      true,
		WithTraceID:        false,
		WithSpanID:         false,
		WithRequestID:      true,
		WithRequestBody:    true,
		WithRequestHeader:  false,
		WithResponseBody:   true,
		WithResponseHeader: false,
		WithXRay:           false,
	})
}

// LoggerWithConfig sets up request logging based on the provided Config.
func LoggerWithConfig(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		attributes := []slog.Attr{}

		// ---------- 1. Base Group ----------

		// user-agent
		if config.WithUserAgent {
			userAgent := c.Get(fiber.HeaderUserAgent)
			attributes = append(attributes, slog.String("user-agent", userAgent))
		}

		// ---------- 2. Trace Group ----------

		ctx := c.UserContext()

		// request_id
		if config.WithRequestID {
			requestID := c.Get(RequestIDHeaderKey)
			if requestID == "" {
				requestID = uuid.New().String()
				c.Set(RequestIDHeaderKey, requestID) // Set for response header
			}
			c.Locals(requestIDCtx, requestID)

			// Add request_id to logx context
			ctx = logx.AppendCtx(ctx, slog.String(RequestIDKey, requestID))
		}

		// trace_id + span_id
		if config.WithTraceID || config.WithSpanID {
			spanCtx := trace.SpanContextFromContext(ctx)

			if config.WithTraceID && spanCtx.HasTraceID() {
				traceID := spanCtx.TraceID().String()
				ctx = logx.AppendCtx(ctx, slog.String("trace_id", traceID))
			}
			if config.WithSpanID && spanCtx.HasSpanID() {
				spanID := spanCtx.SpanID().String()
				ctx = logx.AppendCtx(ctx, slog.String("span_id", spanID))
			}
		}

		c.SetUserContext(ctx)

		// Call the next handler in the middleware chain and capture any error after setting up context but before logging
		err := c.Next()

		// ---------- 3. Request Group ----------

		requestAttributes := []slog.Attr{}

		// request.method
		method := c.Method()
		requestAttributes = append(requestAttributes, slog.String("method", method))

		// request.host
		host := string(c.Context().Host())
		requestAttributes = append(requestAttributes, slog.String("host", host))

		// request.path
		path := c.Path()
		requestAttributes = append(requestAttributes, slog.String("path", path))

		// request.header
		if config.WithRequestHeader {
			kv := []any{}

			for k, v := range c.GetReqHeaders() {
				if _, found := HiddenRequestHeaders[strings.ToLower(k)]; found {
					continue
				}
				kv = append(kv, slog.Any(k, v))
			}
			requestAttributes = append(requestAttributes, slog.Group("header", kv...))
		}

		// request.query
		query := string(c.Request().URI().QueryString())
		requestAttributes = append(requestAttributes, slog.String("query", query))

		// request.param
		params := c.AllParams()
		requestAttributes = append(requestAttributes, slog.Any("params", params))

		// request.body
		if config.WithRequestBody {
			requestBodyBytes := c.Body()
			resquestBodyString := string(requestBodyBytes)
			if len(resquestBodyString) > ResponseBodyMaxSize {
				resquestBodyString = resquestBodyString[:ResponseBodyMaxSize]
			}
			requestAttributes = append(requestAttributes, slog.String("body", resquestBodyString))
		}

		// ---------- 4. Response Group ----------

		responseAttributes := []slog.Attr{}

		// response.status
		status := c.Response().StatusCode()
		responseAttributes = append(responseAttributes, slog.Int("status", status))

		// response.header
		if config.WithResponseHeader {
			kv := []any{}
			for k, v := range c.GetRespHeaders() {
				if _, found := HiddenResponseHeaders[strings.ToLower(k)]; found {
					continue
				}
				kv = append(kv, slog.Any(k, v))
			}

			responseAttributes = append(responseAttributes, slog.Group("header", kv...))
		}

		responseBodyBytes := c.Response().Body()

		// response.body
		if config.WithResponseBody || status >= http.StatusBadRequest {
			responseBodyString := string(responseBodyBytes)
			if len(responseBodyString) > ResponseBodyMaxSize {
				responseBodyString = responseBodyString[:ResponseBodyMaxSize]
			}
			responseAttributes = append(responseAttributes, slog.String("body", responseBodyString))
		}

		// response.length
		responseBodyLength := len(responseBodyBytes)
		responseAttributes = append(responseAttributes, slog.Int("length", responseBodyLength))

		// response.stacktrace
		if status >= http.StatusBadRequest {
			responseAttributes = append(responseAttributes, slog.Any("stacktrace", marshalStack(3)))
		}

		attributes = append(
			attributes,
			slog.Attr{
				Key:   "request",
				Value: slog.GroupValue(requestAttributes...),
			},
			slog.Attr{
				Key:   "response",
				Value: slog.GroupValue(responseAttributes...),
			},
		)

		// custom context values
		if v, ok := c.Locals(customAttributesCtxKey).([]slog.Attr); ok {
			attributes = append(attributes, v...)
		}

		// ---------- Log Message ----------

		level := config.DefaultLevel
		msg := "Incoming request"
		if status >= http.StatusInternalServerError {
			level = config.ServerErrorLevel
			if err != nil {
				msg = err.Error()
			} else {
				msg = http.StatusText(status) // Or a generic server error message
			}
		} else if status >= http.StatusBadRequest {
			level = config.ClientErrorLevel
			if err != nil {
				msg = err.Error()
			} else {
				msg = http.StatusText(status) // Or a generic client error message
			}
		}

		slog.LogAttrs(c.UserContext(), level, msg, attributes...)

		return err // Important to return the error from c.Next()
	}
}

// GetRequestID returns the request identifier.
func GetRequestID(c *fiber.Ctx) string {
	requestIDVal := c.Locals(requestIDCtx)
	if requestIDVal == nil {
		return ""
	}
	if id, ok := requestIDVal.(string); ok {
		return id
	}

	return ""
}

// AddCustomAttributes adds custom attributes to the request context.
func AddCustomAttributes(c *fiber.Ctx, attr slog.Attr) {
	currentAttrsVal := c.Locals(customAttributesCtxKey)
	if currentAttrsVal == nil {
		c.Locals(customAttributesCtxKey, []slog.Attr{attr})
		return
	}

	switch attrs := currentAttrsVal.(type) {
	case []slog.Attr:
		c.Locals(customAttributesCtxKey, append(attrs, attr))
	}
}

type stackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
}

// marshalStack captures the current call stack and parses it into structured frames
func marshalStack(skip int) []stackFrame {
	const depth = 32
	pc := make([]uintptr, depth)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pc[:n])
	var s []stackFrame

	for {
		frame, more := frames.Next()
		if frame.Function == "" {
			break
		}
		s = append(s, stackFrame{
			Func: filepath.Base(frame.Function),
			Source: filepath.Join(
				filepath.Base(filepath.Dir(frame.File)),
				filepath.Base(frame.File),
			),
			Line: frame.Line,
		})

		if !more {
			break
		}
	}
	return s
}
