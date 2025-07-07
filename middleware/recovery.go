package middleware

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Recovery recovers from panics and logs them.
func Recovery(stack bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Simplified request dump for Fiber
				// For more complex dumping, consider c.Request().Header.String() and c.Body() if needed
				requestDump := fmt.Sprintf("%s %s", c.Method(), c.OriginalURL())

				if brokenPipe {
					slog.Error( // Use ErrorCtx to include context attributes like request_id
						c.Path(), // Use c.Path() for Fiber
						slog.Any("error", err),
						slog.String("request", requestDump),
					)
					// If the connection is dead, we can't write a status to it.
					// No explicit abort needed in Fiber for broken pipe, just return.
					// The error won't be sent back to client as connection is broken.
					return
				}

				logArgs := []any{
					slog.Any("error", err),
					slog.String("request", requestDump),
				}
				if stack {
					logArgs = append(logArgs, slog.String("stack", string(debug.Stack())))
				}
				slog.Error("[PANIC RECOVER]", logArgs...)

				// If response has not been sent, send a 500 error.
				// Fiber's default error handler might also catch this panic if not "handled" here.
				// By sending a response, we effectively handle it.
				c.Status(http.StatusInternalServerError)
			}
		}()
		return c.Next() // Call the next handler in the chain.
	}
}
