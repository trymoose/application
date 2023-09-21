package handler

import (
	"io"
	"log/slog"
)

type Handler string

const (
	JSON Handler = "json"
	Text Handler = "text"
)

type HandlerOptions struct {
	Handler Handler `long:"log-handler" choice:"json" choice:"text" default:"text" description:"Handler to use for structured logging." env:"APP_LOGGER_HANDLER"`
}

type _HandlerConstructor func(io.Writer, *slog.HandlerOptions) slog.Handler

func (h *HandlerOptions) Get() _HandlerConstructor {
	switch h.Handler {
	case JSON:
		return func(w io.Writer, o *slog.HandlerOptions) slog.Handler { return slog.NewJSONHandler(w, o) }
	case Text:
		return func(w io.Writer, o *slog.HandlerOptions) slog.Handler { return slog.NewTextHandler(w, o) }
	default:
		return (&HandlerOptions{Handler: Text}).Get()
	}
}
