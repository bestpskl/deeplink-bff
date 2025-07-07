package main

import (
	"deeplink-bff/pkg/logx"
	"log/slog"
	"os"
)

// UserCredentials demonstrates struct with sensitive tags
type UserCredentials struct {
	Username  string `json:"username"`
	Password  string `json:"password" sensitive:"true"`
	APIKey    string `json:"api_key" sensitive:"true"`
	PublicKey string `json:"public_key"`
}

// UserProfile demonstrates nested structs with sensitive data
type UserProfile struct {
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	Age          int             `json:"age"`
	Credentials  UserCredentials `json:"credentials"`
	PhoneNumbers []string        `json:"phone_numbers"`
}

func main() {
	// Initialize logger with configuration
	cfg := logx.Config{
		Environment: "development",
		Source:      "example-app",
	}

	// Configure logger options
	opts := []logx.Option{
		logx.WithLevel(slog.LevelDebug),
		logx.WithSensitiveKeys([]string{"phone_numbers", "custom_secret"}),
	}

	log, err := logx.New(cfg, opts...)
	if err != nil {
		os.Exit(1)
	}

	// Example 1: Basic logging with sensitive data
	log.Info("user login attempt",
		slog.String("username", "john_doe"),
		slog.String("password", "secret123"), // Will be redacted
		slog.String("ip_address", "192.168.1.1"),
	)

	// Example 2: Logging with groups
	log.Info("user profile updated",
		slog.Group("user",
			slog.String("name", "John Doe"),
			slog.String("email", "john@example.com"), // Will be redacted
			slog.Int("age", 30),
		),
	)

	// Example 3: Logging complex structs
	credentials := UserCredentials{
		Username:  "john_doe",
		Password:  "very_secret", // Will be redacted
		APIKey:    "api_key_123", // Will be redacted
		PublicKey: "public_key_xyz",
	}

	profile := UserProfile{
		Name:         "John Doe",
		Email:        "john@example.com", // Will be redacted
		Age:          30,
		Credentials:  credentials,
		PhoneNumbers: []string{"123-456-7890"}, // Will be redacted
	}

	log.Info("user profile created",
		slog.Any("profile", profile),
	)

	// Example 4: Logging JSON strings
	jsonData := `{
			"user": {
					"email": "john@example.com",
					"password": "secret123",
					"preferences": {
							"theme": "dark",
							"api_key": "key123"
					}
			}
	}`

	log.Info("received user data",
		slog.String("data", jsonData),
	)

	// Example 5: Debug logging
	log.Debug("debug information",
		slog.String("api_key", "debug_key"), // Will be redacted
		slog.String("request_id", "req123"),
	)

	// Example 6: Error logging with sensitive data
	log.Error("authentication failed",
		slog.String("username", "john_doe"),
		slog.String("password", "wrong_pass"), // Will be redacted
		slog.String("error", "invalid credentials"),
	)

	// Example 7: Logging with custom attributes
	log.Info("custom operation",
		slog.Group("metadata",
			slog.String("custom_secret", "hidden_value"), // Will be redacted
			slog.String("visible_field", "normal_value"),
		),
	)
}
