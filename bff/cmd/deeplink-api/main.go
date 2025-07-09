package main

import (
	"deeplink-bff/bff/config"
	"deeplink-bff/bff/docs"
	deeplink_client "deeplink-bff/bff/internal/adapters/client"
	deeplink_handler "deeplink-bff/bff/internal/adapters/handler/deeplink"
	deeplink_service "deeplink-bff/bff/internal/core/services/deeplink"
	"deeplink-bff/middleware"
	"deeplink-bff/pkg/logx"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @securityDefinitions.apikey	Authorization
// @in							header
// @name						Authorization
// @description				Please input prefix "Bearer " and your access token.
// @securityDefinitions.apikey	X-OPENAPI-JWT
// @in							header
// @name						X-OPENAPI-JWT
// @description				Please input your customer id. (works only in dev environment)
func newRouters(
	deeplinkHandler *deeplink_handler.Handler,
) *fiber.App {
	appConfig := fiber.Config{
		// Fiber's default error handler is quite good.
		// ErrorHandler: func(c *fiber.Ctx, err error) error { ... }
	}
	if config.Get().Environment != "dev" {
		// Mimic Gin's ReleaseMode effects for non-dev environments
		appConfig.DisableStartupMessage = true // Suppress Fiber's startup banner
	}

	app := fiber.New(appConfig)
	// app.Use(middleware.Error()) // Ensure middleware.Error() is Fiber compatible if used

	if config.Get().Environment != "staging" && config.Get().Environment != "prod" {
		app.Get("/docs/*", initSwagger()) // Use "/*" for path parameters in Fiber
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{ // Use fiber.Map and return error
			"message": "ok",
		})
	})

	apiGroup := app.Group("/api")
	v1 := apiGroup.Group("/v1")
	v1.Use(
		// middleware.Auth(), // Ensure this middleware is Fiber compatible: func(c *fiber.Ctx) error
		middleware.Logger(),
		middleware.Recovery(true),
	)

	dashboardGroup := v1.Group("/deeplink")
	{
		dashboardGroup.Get("", deeplinkHandler.GetDeeplinkList)
		// Ensure deeplinkHandler.GetDeeplinkList signature is: func(c *fiber.Ctx) error
		dashboardGroup.Get("/:id", deeplinkHandler.GetDeeplink)
	}
	return app
}

func initSwagger() fiber.Handler { // Return type changed to fiber.Handler
	docs.SwaggerInfo.Title = "deeplink API"
	docs.SwaggerInfo.Description = "APIs for providing deeplink data."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.BasePath = "/api"

	if config.Get().Environment != "dev" {
		docs.SwaggerInfo.Schemes = []string{"https"}
		// docs.SwaggerInfo.BasePath = "/deeplink/api"
	}

	return fiberSwagger.WrapHandler // fiberSwagger.WrapHandler is the fiber.Handler
}

func main() {
	config.Load()

	cfgLog := logx.Config{
		Environment: config.Get().Environment,
		Source:      "deeplink-bff",
	}

	// local file-based logging
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Open the log file
	logFilePath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	writer := io.MultiWriter(os.Stdout, logFile)

	logger, _ := logx.New(
		cfgLog,
		logx.WithAddSource(false),
		logx.WithLevel(slog.LevelDebug),
		logx.WithWriter(writer),
	)

	slog.SetDefault(logger)

	deeplinkClient := deeplink_client.NewDeepLinkClient("http://localhost:3000")
	deeplinkService := deeplink_service.NewDeeplinkService(deeplinkClient)
	deeplinkHandler := deeplink_handler.NewHandler(deeplinkService)
	app := newRouters(deeplinkHandler)

	addr := fmt.Sprintf("%s:%d", config.Get().App.Host, config.Get().App.Port)

	// Start server in a goroutine to allow for graceful shutdown
	go func() {
		slog.Info("Starting server", slog.String("addr", addr))
		// app.Listen blocks until the server is shut down or an error occurs.
		// If shutdown is graceful (via app.Shutdown()), Listen returns nil.
		if err := app.Listen(addr); err != nil {
			slog.Error("Server Listen error", slog.Any("error", err))
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received

	slog.Info("Shutting down server...")

	// Attempt to gracefully shut down the server with a timeout.
	shutdownTimeout := 1 * time.Minute // Define a timeout for shutdown
	if err := app.ShutdownWithTimeout(shutdownTimeout); err != nil {
		slog.Error("Server forced to shutdown:", slog.Any("error", err))
	}

	slog.Info("Server exiting")
}
