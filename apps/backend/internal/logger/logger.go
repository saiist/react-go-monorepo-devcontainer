package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type contextKey string

const loggerKey contextKey = "logger"

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
}

// Config represents logger configuration
type Config struct {
	Level      string
	Format     string // "json" or "text"
	TimeFormat string
}

// New creates a new logger instance
func New(cfg Config) *Logger {
	var handler slog.Handler

	level := parseLevel(cfg.Level)
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// GCP Cloud Loggingに適した形式に変換
			switch a.Key {
			case slog.LevelKey:
				// Cloud Loggingのseverityフィールドに変換
				return slog.Attr{
					Key:   "severity",
					Value: slog.StringValue(gcpSeverity(a.Value.Any().(slog.Level))),
				}
			case slog.TimeKey:
				// RFC3339形式のタイムスタンプ
				return slog.Attr{
					Key:   "timestamp",
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339Nano)),
				}
			case slog.SourceKey:
				// ソースの場所を短縮
				if src, ok := a.Value.Any().(*slog.Source); ok {
					return slog.Attr{
						Key: "source",
						Value: slog.GroupValue(
							slog.String("file", src.File),
							slog.Int("line", src.Line),
							slog.String("function", src.Function),
						),
					}
				}
			}
			return a
		},
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithContext returns a logger from context or creates a new one
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	return l
}

// ToContext adds logger to context
func (l *Logger) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext extracts logger from context
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	// デフォルトロガーを返す
	return New(Config{Level: "info", Format: "json"})
}

// WithRequest adds request-specific fields
func (l *Logger) WithRequest(requestID, traceID, method, path, remoteAddr string) *Logger {
	return &Logger{
		Logger: l.With(
			slog.String("request_id", requestID),
			slog.String("trace_id", traceID),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("remote_addr", remoteAddr),
		),
	}
}

// WithUser adds user context
func (l *Logger) WithUser(userID, email string) *Logger {
	return &Logger{
		Logger: l.With(
			slog.String("user_id", userID),
			slog.String("user_email", email),
		),
	}
}

// WithError adds error with stack trace
func (l *Logger) WithError(err error) *Logger {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	
	var stack []slog.Attr
	for {
		frame, more := frames.Next()
		stack = append(stack, slog.Group("frame",
			slog.String("file", frame.File),
			slog.Int("line", frame.Line),
			slog.String("function", frame.Function),
		))
		if !more {
			break
		}
	}

	return &Logger{
		Logger: l.With(
			slog.String("error", err.Error()),
			slog.Any("stack", stack),
		),
	}
}

// Helper functions

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func gcpSeverity(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARNING"
	case slog.LevelError:
		return "ERROR"
	default:
		return "DEFAULT"
	}
}