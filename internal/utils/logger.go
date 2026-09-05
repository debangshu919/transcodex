package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/debangshu919/transcodex/internal/config"
)

// filteredHandler forwards log records to an underlying slog.Handler only when match returns true.
type filteredHandler struct {
	handler slog.Handler
	match   func(slog.Level) bool
}

func (h *filteredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.match(level) && h.handler.Enabled(ctx, level)
}

func (h *filteredHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.match(r.Level) {
		return nil
	}
	return h.handler.Handle(ctx, r)
}

func (h *filteredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &filteredHandler{
		handler: h.handler.WithAttrs(attrs),
		match:   h.match,
	}
}

func (h *filteredHandler) WithGroup(name string) slog.Handler {
	return &filteredHandler{
		handler: h.handler.WithGroup(name),
		match:   h.match,
	}
}

// multiHandler broadcasts log records to multiple slog.Handlers.
type multiHandler struct {
	handlers []slog.Handler
	closers  []io.Closer
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{
		handlers: handlers,
		closers:  m.closers,
	}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{
		handlers: handlers,
		closers:  m.closers,
	}
}

func (m *multiHandler) Close() error {
	var errs []error
	for _, c := range m.closers {
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// parseLogLevel converts a level string into a slog.Level.
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewLogger initializes and returns a *slog.Logger configured for stdout and JSON file logs.
func NewLogger(cfg config.Config, logsDir string) *slog.Logger {
	return newLoggerWithWriter(cfg, logsDir, os.Stdout)
}

func newLoggerWithWriter(cfg config.Config, logsDir string, stdout io.Writer) *slog.Logger {
	if strings.TrimSpace(logsDir) == "" {
		logsDir = "logs"
	}
	logsDir = filepath.Clean(logsDir)

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logs directory %s: %v\n", logsDir, err)
	}

	stdoutLevel := parseLogLevel(cfg.LogLevel)
	stdoutHandler := slog.NewTextHandler(stdout, &slog.HandlerOptions{
		Level:     stdoutLevel,
		AddSource: false,
	})

	handlers := []slog.Handler{stdoutHandler}
	var closers []io.Closer

	// app.log: always logs info logs (JSON, without source)
	appLogPath := filepath.Join(logsDir, "app.log")
	appFile, err := os.OpenFile(appLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open app log file %s: %v\n", appLogPath, err)
	} else {
		closers = append(closers, appFile)
		appJSONHandler := slog.NewJSONHandler(appFile, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		})
		handlers = append(handlers, &filteredHandler{
			handler: appJSONHandler,
			match: func(l slog.Level) bool {
				return l >= slog.LevelInfo && l < slog.LevelError
			},
		})
	}

	// error.log: always logs error logs (JSON, with source)
	errorLogPath := filepath.Join(logsDir, "error.log")
	errorFile, err := os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open error log file %s: %v\n", errorLogPath, err)
	} else {
		closers = append(closers, errorFile)
		errorJSONHandler := slog.NewJSONHandler(errorFile, &slog.HandlerOptions{
			Level:     slog.LevelError,
			AddSource: true,
		})
		handlers = append(handlers, &filteredHandler{
			handler: errorJSONHandler,
			match: func(l slog.Level) bool {
				return l >= slog.LevelError
			},
		})
	}

	// debug.log: only created when log_level is set to debug (JSON, with source)
	if strings.EqualFold(strings.TrimSpace(cfg.LogLevel), "debug") {
		debugLogPath := filepath.Join(logsDir, "debug.log")
		debugFile, err := os.OpenFile(debugLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open debug log file %s: %v\n", debugLogPath, err)
		} else {
			closers = append(closers, debugFile)
			debugJSONHandler := slog.NewJSONHandler(debugFile, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			})
			handlers = append(handlers, &filteredHandler{
				handler: debugJSONHandler,
				match: func(l slog.Level) bool {
					return l < slog.LevelInfo
				},
			})
		}
	}

	mh := &multiHandler{
		handlers: handlers,
		closers:  closers,
	}

	return slog.New(mh)
}

// closeLogger closes any underlying open log files.
func closeLogger(l *slog.Logger) error {
	if l == nil {
		return nil
	}
	if closer, ok := l.Handler().(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
