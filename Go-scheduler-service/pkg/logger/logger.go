package logger

import (
	"github.com/rs/zerolog"
)

type ZeroLogLogger struct {
	logger *zerolog.Logger
}

func New(level string, logger *zerolog.Logger) (*ZeroLogLogger, error) {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	lvlLogger := logger.Level(parsedLevel)
	logger = &lvlLogger

	return &ZeroLogLogger{
		logger: logger,
	}, nil
}

func (l *ZeroLogLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l *ZeroLogLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *ZeroLogLogger) Warning(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l *ZeroLogLogger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

func (l *ZeroLogLogger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}
