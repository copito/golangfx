package modules

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
)

type LoggerParams struct {
	fx.In
}

type LoggerResult struct {
	fx.Out

	Logger *slog.Logger
}

// func NewLogger(p LoggerParams) (LoggerResult, error) {
// 	loggerWithColor := os.Getenv("LOGGER_WITH_COLOR")
// 	loggerLevel := os.Getenv("LOGGER_LEVEL")

// 	var withColor bool
// 	if strings.ToLower(loggerWithColor) == "true" {
// 		withColor = true
// 	} else {
// 		withColor = false
// 	}

// 	var level slog.Level
// 	if strings.ToLower(loggerLevel) == "debug" {
// 		level = slog.LevelDebug
// 	} else if strings.ToLower(loggerLevel) == "info" {
// 		level = slog.LevelInfo
// 	} else if strings.ToLower(loggerLevel) == "warn" {
// 		level = slog.LevelWarn
// 	} else if strings.ToLower(loggerLevel) == "error" {
// 		level = slog.LevelError
// 	} else {
// 		level = slog.LevelInfo
// 	}

// 	handler := logger.NewCustomLogHandler(withColor, level)
// 	logger := slog.New(handler)

// 	return LoggerResult{
// 		Logger: logger,
// 	}, nil
// }

func NewLogger(p LoggerParams) (LoggerResult, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("setting up logging (with slog)...")

	return LoggerResult{
		Logger: logger,
	}, nil
}

var LoggerModule = fx.Provide(NewLogger)
