package modules

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/copito/runner/src/internal/entities"
	"go.uber.org/fx"
)

type LoggerParams struct {
	fx.In
	Config *entities.Config
}

type LoggerResult struct {
	fx.Out

	Logger *slog.Logger
}

func ReplaceSourceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		if s, ok := a.Value.Any().(*slog.Source); ok {
			filename := filepath.Base(s.File)
			return slog.Any(a.Key, filename)
		}
	}
	return a
}

func NewLogger(p LoggerParams) (LoggerResult, error) {
	var level slog.Level
	switch p.Config.Logger.Level {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: ReplaceSourceAttr,
		Level:       level,
	}

	var loggerHandler slog.Handler
	switch p.Config.Logger.Type {
	case "JSON":
		loggerHandler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	case "TEXT":
		loggerHandler = slog.NewTextHandler(os.Stdout, handlerOptions)
	default:
		loggerHandler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	}

	logger := slog.New(loggerHandler)
	logger.Info("setting up logging module (with slog)...")

	return LoggerResult{
		Logger: logger,
	}, nil
}

var LoggerModule = fx.Provide(NewLogger)
