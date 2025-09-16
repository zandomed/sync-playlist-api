package logger

import (
	"fmt"

	"github.com/zandomed/sync-playlist-api/internal/config"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	var logger *zap.Logger
	var err error

	cfg := config.Get()
	// Configurar logger basado en ambiente

	if cfg.IsDevelopment() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	return &Logger{Logger: logger}
}
