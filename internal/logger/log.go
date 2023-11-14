package logger

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
	log.Trace().Msg("Zerolog initialized.")
}
