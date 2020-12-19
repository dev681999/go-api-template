package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger is default logger
var Logger zerolog.Logger

// Setup sets up logger
func Setup() zerolog.Logger {
	logger := zerolog.New(os.Stdout)
	logger = logger.With().Timestamp().Caller().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Logger = logger
	Logger = logger

	return logger
}
