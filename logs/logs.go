package logs

import (
	"log"
	"os"
	"time"

	"github.com/daqiancode/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Log zerolog.Logger

func init() {
	level, err := zerolog.ParseLevel(env.Get("LOG_LEVEL", "DEBUG"))
	if err != nil {
		log.Fatal("Error loading LOG_LEVEL", err)
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DisableSampling(true)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true}
	Log = zerolog.New(output).With().Timestamp().Logger()
}
