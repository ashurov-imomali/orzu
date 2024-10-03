package logger

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"os"
)

type logger struct {
	z zerolog.Logger
}

func New() Logger {
	return logger{z: zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(3).Logger()}
}

func (l logger) Printf(format string, a ...interface{}) {
	l.z.Printf(format, a...)
}

func (l logger) Println(a ...interface{}) {
	for _, i := range a {
		if marshal, err := json.Marshal(&i); err == nil {
			l.z.Print(string(marshal))
		}
	}
}

func (l logger) Error(a ...interface{}) {
	l.z.Error().Msgf("%v", a)
}

func (l logger) Errorf(format string, a ...interface{}) {
	l.z.Error().Msgf(format, a)
}

func (l logger) KVLog(k string, v interface{}) {
	l.z.Info().Interface(k, v).Send()
}

func (l logger) Fatal(a ...interface{}) {
	l.z.Fatal().Msgf("%v", a)
}

func (l logger) Warn(a ...interface{}) {
	l.z.Warn().Msgf("%v", a)
}

func (l logger) Info(a ...interface{}) {
	l.z.Info().Msgf("%v", a)
}
