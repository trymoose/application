package flags

import (
	"context"
	"log/slog"
)

type _ContextKey string

const (
	_ContextKeyArgs   _ContextKey = "args"
	_ContextKeyExit   _ContextKey = "exit"
	_ContextKeyLogger _ContextKey = "logger"
)

type ModCtx interface {
	ModCtx(context.Context) (context.Context, error)
}

func ContextArgs(ctx context.Context) []string {
	return ctx.Value(_ContextKeyArgs).([]string)
}

func ContextExit[I ~int](ctx context.Context, code I) {
	ctx.Value(_ContextKeyExit).(_Exiter)(ctx, int(code))
}

func ContextLogger(ctx context.Context) *slog.Logger {
	return ctx.Value(_ContextKeyLogger).(*slog.Logger)
}

func ContextLoggerWith(ctx context.Context, args ...any) context.Context {
	return context.WithValue(ctx, _ContextKeyLogger, ContextLogger(ctx).With(args...))
}
