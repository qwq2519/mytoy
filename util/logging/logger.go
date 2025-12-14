package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"mytoy/config"
)

// NewLogger 根据配置创建 slog.Logger，并返回复用的 writer（用于 Gin 等标准输出）。
func NewLogger(cfg config.LoggingConfig) (*slog.Logger, io.Writer, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, nil, err
	}

	if err := os.MkdirAll(filepath.Dir(cfg.File), 0o755); err != nil {
		return nil, nil, fmt.Errorf("ensure log dir: %w", err)
	}

	file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	writers := []io.Writer{file}
	if cfg.EnableConsole {
		writers = append(writers, os.Stdout)
	}
	multi := io.MultiWriter(writers...)

	handler := slog.NewJSONHandler(multi, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)

	return logger, multi, nil
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unsupported logging.level: %s", level)
	}
}
