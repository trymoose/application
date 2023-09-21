//go:build !debug

package level

import (
	"log/slog"
)

const Default = Info

type LevelOptions struct {
	Level Level `long:"log-level" choice:"info" choice:"error" choice:"warn" choice:"debug" description:"Log level to log at." default:"info" env:"APP_LOGGER_LEVEL"`
}

func (o LevelOptions) Get() slog.Leveler { return o.Level.Level() }
