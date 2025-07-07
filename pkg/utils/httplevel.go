package utils

import (
	"log/slog"
	"net/http"
)

func HttpStatusCodeToLogLevel(code int) slog.Level {
	lv := slog.LevelInfo
	switch {
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		lv = slog.LevelWarn
	case code >= http.StatusInternalServerError:
		lv = slog.LevelError
	default:
		lv = slog.LevelInfo
	}
	return lv
}
