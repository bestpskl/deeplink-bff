package logx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type logOutput struct {
	Level   string                 `json:"level"`
	Msg     string                 `json:"msg"`
	User    map[string]interface{} `json:"user"`
	Runtime map[string]interface{} `json:"runtime"`
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		opts    []Option
		wantErr bool
		check   func(*testing.T, *slog.Logger)
	}{
		{
			name: "basic configuration",
			cfg: Config{
				Environment: "test",
				Source:      "test-app",
			},
			check: func(t *testing.T, logger *slog.Logger) {
				assert.NotNil(t, logger)
			},
		},
		{
			name: "with custom redaction keys",
			cfg: Config{
				Environment: "test",
				Source:      "test-app",
			},
			opts: []Option{
				WithSensitiveKeys([]string{"custom_secret"}),
			},
			check: func(t *testing.T, logger *slog.Logger) {
				assert.NotNil(t, logger)
				// Verify custom key is in redactKeys
				_, exists := sensitivekeys["custom_secret"]
				assert.True(t, exists)
			},
		},
		{
			name: "with debug mode",
			cfg: Config{
				Environment: "test",
				Source:      "test-app",
			},
			opts: []Option{
				WithDebugMode(true),
			},
			check: func(t *testing.T, logger *slog.Logger) {
				assert.NotNil(t, logger)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.cfg, tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.check(t, logger)
		})
	}
}

func TestRedactingHandler(t *testing.T) {
	tests := []struct {
		name     string
		input    func(*slog.Logger)
		wantLogs func(*testing.T, logOutput)
	}{
		{
			name: "redacts sensitive field",
			input: func(log *slog.Logger) {
				log.Info("test message",
					slog.Group("user",
						slog.String("email", "test@example.com"),
						slog.String("name", "John"),
					),
				)
			},
			wantLogs: func(t *testing.T, output logOutput) {
				assert.Equal(t, "INFO", output.Level)
				assert.Equal(t, "test message", output.Msg)
				userMap := output.User
				assert.Equal(t, DefaultRedactMessage, userMap["email"])
				assert.Equal(t, "John", userMap["name"])
			},
		},
		{
			name: "redacts nested sensitive fields",
			input: func(log *slog.Logger) {
				log.Info("test message",
					slog.Group("user",
						slog.String("credentials", `{"email":"test@example.com","api_key":"secret"}`),
					),
				)
			},
			wantLogs: func(t *testing.T, output logOutput) {
				assert.Equal(t, "INFO", output.Level)
				userMap := output.User
				credentials, ok := userMap["credentials"].(string)
				require.True(t, ok, "credentials should be a string")

				// Parse the credentials JSON string
				var credMap map[string]interface{}
				err := json.Unmarshal([]byte(credentials), &credMap)
				require.NoError(t, err, "should be able to parse credentials JSON")

				assert.Equal(t, DefaultRedactMessage, credMap["email"])
				assert.Equal(t, DefaultRedactMessage, credMap["api_key"])
			},
		},
		{
			name: "handles non-json string field",
			input: func(log *slog.Logger) {
				log.Info("test message",
					slog.Group("user",
						slog.String("email", "test@example.com"),
						slog.String("regular_field", "normal value"),
					),
				)
			},
			wantLogs: func(t *testing.T, output logOutput) {
				assert.Equal(t, "INFO", output.Level)
				userMap := output.User
				assert.Equal(t, DefaultRedactMessage, userMap["email"])
				assert.Equal(t, "normal value", userMap["regular_field"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture log output
			buf := &bytes.Buffer{}

			// Create logger with buffer
			opts := &slog.HandlerOptions{Level: slog.LevelDebug}
			handler := slog.NewJSONHandler(buf, opts)

			sensitivekeys = make(map[string]struct{})
			// Set sensitive keys from blacklist by default
			for key := range blackList {
				sensitivekeys[key] = struct{}{}
			}

			censoringHandler := &censoringHandler{
				handler:       handler,
				withDebug:     false,
				sensitiveKeys: sensitivekeys,
			}

			logger := slog.New(censoringHandler)

			// Execute test
			tt.input(logger)

			// Parse output
			var output logOutput
			err := json.Unmarshal(buf.Bytes(), &output)
			require.NoError(t, err)

			// Verify output
			tt.wantLogs(t, output)
		})
	}
}

func TestOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   Option
		validate func(*testing.T, *optionsConfig)
	}{
		{
			name:   "WithDebugMode",
			option: WithDebugMode(true),
			validate: func(t *testing.T, opts *optionsConfig) {
				assert.True(t, opts.WithDebug)
			},
		},
		{
			name:   "WithLevel",
			option: WithLevel(slog.LevelWarn),
			validate: func(t *testing.T, opts *optionsConfig) {
				assert.Equal(t, slog.LevelWarn, opts.Level)
			},
		},
		{
			name:   "WithSensitiveKeys",
			option: WithSensitiveKeys([]string{"custom_key"}),
			validate: func(t *testing.T, opts *optionsConfig) {
				assert.Contains(t, opts.SensitiveKeys, "custom_key")
			},
		},
		{
			name:   "WithDefaultRedactMessage",
			option: WithDefaultRedactMessage("***HIDDEN***"),
			validate: func(t *testing.T, opts *optionsConfig) {
				assert.Equal(t, "***HIDDEN***", opts.DefaultRedactMessage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &optionsConfig{}
			err := tt.option(opts)
			assert.NoError(t, err)
			tt.validate(t, opts)
		})
	}
}
