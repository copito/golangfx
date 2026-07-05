package log

import (
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/internal/modules/config"
)

type Params struct {
	fx.In
	ConfigProvider config.ConfigProvider
}

type Result struct {
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

func NewLogger(p Params) (Result, error) {
	config := p.ConfigProvider.Get()

	var level slog.Level
	switch config.Logger.Level {
	case entities.LoggerLevelDEBUG:
		level = slog.LevelDebug
	case entities.LoggerLevelINFO:
		level = slog.LevelInfo
	case entities.LoggerLevelWARN:
		level = slog.LevelWarn
	case entities.LoggerLevelERROR:
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
	switch config.Logger.Type {
	case entities.LoggerTypeJSON:
		loggerHandler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	case entities.LoggerTypeTEXT:
		loggerHandler = slog.NewTextHandler(os.Stdout, handlerOptions)
	default:
		loggerHandler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	}

	logger := slog.New(loggerHandler)
	logger.Info("setting up logging module (with slog)...")

	return Result{
		Logger: logger,
	}, nil
}

var LoggerModule = fx.Provide(NewLogger)
