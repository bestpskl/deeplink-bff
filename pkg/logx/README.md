# Logging Package

## Overview

This package provides a structured logging utility for Go applications using `slog`. It includes features such as log level management, source tracking, and sensitive data masking to enhance security and debugging.

## Features

- **Configurable Log Levels**: Supports different log levels (Debug, Info, Warn, Error).
- **JSON Logging**: Uses structured JSON format for logs.
- **Sensitive Data Censoring**: Masks specified sensitive keys in logs.
- **Source Tracking**: Optionally includes source file and line numbers.

## Installation

```sh
# Add the package to your Go module
 go get github.com/deeplink-bff/pkg/logging
```

## Usage

### Basic Logger Initialization

```go
package main

import (
 "log/slog"
 "deeplink-bff/pkg/logging"
)

func main() {
 logger, err := logging.NewLogger()
 if err != nil {
  panic(err)
 }

 slog.SetDefault(logger)
 slog.Info("Application started")
}
```

### Configuring the Logger

```go
logger, _ := NewLogger(
  WithEnvironment("development"),
  WithAddSource(false),
  WithSource("example-service"),
  WithLevel(slog.LevelDebug),
  WithSensitiveKeys([]string{"full_name", "email"}),
  WithDebug(true),
 )

slog.SetDefault(logger)
```

#### Configuration Options

- `WithLevel(level slog.Level) Option`
  - Sets the log level to control the verbosity of logs.
  - Available levels: `slog.LevelDebug`, `slog.LevelInfo`, `slog.LevelWarn`, `slog.LevelError`.

- `WithAddSource(addSource bool) Option`
  - Enables or disables the inclusion of source information (e.g., filename, line number) in logs.

- `WithSensitiveKeys(keys []string) Option`
  - Defines keys that should be treated as sensitive and redacted in logs.
  - Example: `[]string{"password", "api_key"}`

- `WithDebug(withDebug bool) Option`
  - Enables or disables debug mode.
  - When enabled, additional debugging information may be included in log output.

### Logging Messages

```go
logger.Info("User login successful", slog.String("user", "john_doe"))
// Output: {"time":"2025-02-04T10:14:51.605155+07:00","level":"INFO","msg":"User login successful","source":"example-service","environment":"development","user":"john_doe"}

 err := errors.New("Failed to process request")
 slog.ErrorContext(context.TODO(), "test", slog.Any("error", err))
  // Output: {"time":"2025-02-04T10:15:17.959272+07:00","level":"ERROR","msg":"test","source":"example-service","environment":"development","error":{"msg":"Failed to process request"}}

```

### Masking Sensitive Data

Sensitive fields defined in `SensitiveKeys` will be masked automatically:

```go
logger.Info("User registered", slog.String("password", "supersecret"))
// Output: { "password": "*" }
```

## Advanced Features

### Censoring Handler

The package includes a custom `CensoringHandler` that automatically censors sensitive fields based on configuration.

## License

Krungthai
