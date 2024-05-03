package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// If in development mode use pretty printed output, else log as json, discard logs when testing
func getLoggerStyle(env string) io.Writer {
	if env == "test" {
		return io.Discard
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05.000Z07:00"}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	return output
}

func New(env string) *zerolog.Logger {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	logger := zerolog.New(getLoggerStyle(env)).With().Timestamp().Logger()

	if env == "prod" {
		DisableDebugLogs()
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	return &logger
}

func DisableDebugLogs() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
