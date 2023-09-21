package level

import "log/slog"

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
)

func (l Level) Level() slog.Leveler {
	switch l {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Error:
		return slog.LevelError
	case Warn:
		return slog.LevelWarn
	default:
		return Default.Level()
	}
}
