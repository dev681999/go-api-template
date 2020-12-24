package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger is default logger
var Logger zerolog.Logger

// Env for logging
type Env uint

// Envs
const (
	EnvLocal Env = iota
	EnvTest
	EnvStage
	EnvProd
)

func setLoggerLevel(logger zerolog.Logger, env Env) zerolog.Logger {
	switch env {
	case EnvProd:
		return logger.Level(zerolog.InfoLevel)
	}

	return logger
}

// Setup sets up logger
func Setup(env Env) zerolog.Logger {
	logger := zerolog.New(os.Stdout)
	logger = logger.With().Timestamp().Caller().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger = setLoggerLevel(logger, env)

	log.Logger = logger
	Logger = logger

	return logger
}
