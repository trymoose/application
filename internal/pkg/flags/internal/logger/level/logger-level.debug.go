//go:build debug

package level

import (
	"log/slog"
)

const Default = Debug

type LevelOptions struct {
	Level Level `long:"log-level" choice:"info" choice:"error" choice:"warn" choice:"debug" description:"Log level to log at." default:"debug" env:"APP_LOGGER_LEVEL"`
}

func (o LevelOptions) Get() slog.Leveler { return o.Level.Level() }
