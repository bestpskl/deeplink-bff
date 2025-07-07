
# Gin middleware

## 🚀 Install

```sh
go get deeplink-bff/middleware
```

**Compatibility**: go >= 1.21

No breaking changes will be made to exported APIs before v2.0.0.

## 💡 Usage

### Handler options

```go
type Config struct {
 DefaultLevel     slog.Level
 ClientErrorLevel slog.Level
 ServerErrorLevel slog.Level

 WithUserAgent      bool
 WithRequestID      bool
 WithRequestBody    bool
 WithRequestHeader  bool
 WithResponseBody   bool
 WithResponseHeader bool
 WithSpanID         bool
 WithTraceID        bool

 Filters []Filter
}
```

Attributes will be injected in log payload.

Other global parameters:

```go
sloggin.TraceIDKey = "trace_id"
sloggin.SpanIDKey = "span_id"
sloggin.RequestBodyMaxSize  = 64 * 1024 // 64KB
sloggin.ResponseBodyMaxSize = 64 * 1024 // 64KB
sloggin.HiddenRequestHeaders = map[string]struct{}{ ... }
sloggin.HiddenResponseHeaders = map[string]struct{}{ ... }
sloggin.RequestIDHeaderKey = "X-Request-Id"
```

### Minimal

```go
import (
 "github.com/gin-gonic/gin"
 "deeplink-bff/middleware"
 "log/slog"
)

// Create a slog logger, which:
//   - Logs to stdout.
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

router := gin.New()

// Add the sloggin middleware to all routes.
// The middleware will log all requests attributes.
router.Use(middleware.Logger())
router.Use(gin.Recovery())

// Example pong request.
router.GET("/pong", func(c *gin.Context) {
    c.String(http.StatusOK, "pong")
})

router.Run(":1234")

// output:
// time=2023-10-15T20:32:58.926+02:00 level=INFO msg="Incoming request" env=production request.time=2023-10-15T20:32:58.626+02:00 request.method=GET request.path=/ request.query="" request.route="" request.ip=127.0.0.1:63932 request.length=0 response.time=2023-10-15T20:32:58.926+02:00 response.latency_ms=100ms response.status=200 response.length=7 id=""
```

### OTEL

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

config := middleware.Config{
 WithSpanID:  true,
 WithTraceID: true,
}
router := gin.New()
router.Use(middleware.LoggerWithConfig(logger, config))
```

### Custom log levels

```go
logger := slog.New()

config := middleware.Config{
 DefaultLevel:     slog.LevelInfo,
 ClientErrorLevel: slog.LevelWarn,
 ServerErrorLevel: slog.LevelError,
}

router := gin.New()
router.Use(sloggin.NewWithConfig(logger, config))
```

### Verbose

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

config := middleware.Config{
 WithRequestBody: true,
 WithResponseBody: true,
 WithRequestHeader: true,
 WithResponseHeader: true,
}

router := gin.New()
router.Use(middleware.LoggerWithConfig(logger, config))
```

### Add logger to a single route

```go
logger := slog.New()

router := gin.New()
router.Use(gin.Recovery())

// Example pong request.
// Add the sloggin middleware to a single routes.
router.GET("/pong", middleware.Logger(), func(c *gin.Context) {
    c.String(http.StatusOK, "pong")
})

router.Run(":1234")
```

### Adding custom attributes

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).
    With("environment", "production").
    With("server", "gin/1.9.0").
    With("server_start_time", time.Now()).
    With("gin_mode", gin.EnvGinMode)

router := gin.New()

// Add the sloggin middleware to all routes.
// The middleware will log all requests attributes.
router.Use(middleware.Logger())
router.Use(gin.Recovery())

// Example pong request.
router.GET("/pong", func(c *gin.Context) {
 // Add an attribute to a single log entry.
 sloggin.AddCustomAttributes(c, slog.String("foo", "bar"))
    c.String(http.StatusOK, "pong")
})

router.Run(":1234")

// output:
// time=2023-10-15T20:32:58.926+02:00 level=INFO msg="Incoming request" environment=production server=gin/1.9.0 gin_mode=release request.time=2023-10-15T20:32:58.626+02:00 request.method=GET request.path=/ request.query="" request.route="" request.ip=127.0.0.1:63932 request.length=0 response.time=2023-10-15T20:32:58.926+02:00 response.latency_ms=100ms response.status=200 response.length=7 id="" foo=bar
```

### JSON output

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

router := gin.New()

// Add the sloggin middleware to all routes.
// The middleware will log all requests attributes.
router.Use(middleware.Logger())
router.Use(gin.Recovery())

// Example pong request.
router.GET("/pong", func(c *gin.Context) {
    c.String(http.StatusOK, "pong")
})

router.Run(":1234")

// output:
// {"time":"2023-10-15T20:32:58.926+02:00","level":"INFO","msg":"Incoming request","gin_mode":"GIN_MODE","env":"production","http":{"request":{"time":"2023-10-15T20:32:58.626+02:00","method":"GET","path":"/","query":"","route":"","ip":"127.0.0.1:55296","length":0},"response":{"time":"2023-10-15T20:32:58.926+02:00","latency":100000,"status":200,"length":7},"id":""}}
```
