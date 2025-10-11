package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func Init(environment string) {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var output io.Writer = os.Stdout

	if environment == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Str("service", "data-lake").
		Logger()

	log.Logger = Logger
}

func Get() zerolog.Logger {
	return Logger
}
