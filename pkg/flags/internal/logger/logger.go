package logger

import (
	"context"
	"github.com/trymoose/application/pkg/flags/internal/logger/handler"
	"github.com/trymoose/application/pkg/flags/internal/logger/level"
	"github.com/trymoose/application/pkg/flags/internal/logger/output"
	"github.com/trymoose/application/pkg/flags/internal/logger/source"
	"io"
	"log/slog"
)

const (
	Name  = "logger"
	Short = "Logger options."
	Long  = "Arguments that change the behavior of the logger."
)

type Logger struct {
	level.LevelOptions
	source.SourceOptions
	output.OutputOptions
	handler.HandlerOptions

	_CtxKey   any
	_Close    io.Closer
	_LevelVar *slog.LevelVar
}

func New(ctxKey any, logLevel *slog.LevelVar) *Logger {
	l := &Logger{_CtxKey: ctxKey, _LevelVar: logLevel}
	l._LevelVar.Set(slog.LevelInfo)
	l._SetNoopCloser()
	return l
}

func (l *Logger) _SetNoopCloser() {
	l._Close = io.NopCloser(io.MultiReader())
}

func (l *Logger) Close() (err error) {
	err = l._Close.Close()
	l._SetNoopCloser()
	return
}

func (l *Logger) ModCtx(ctx context.Context) (context.Context, error) {
	f, err := l.OutputOptions.Writer()
	if err != nil {
		return nil, err
	}
	l._Close = f

	l._LevelVar.Set(l.LevelOptions.Get().Level())
	return context.WithValue(ctx, l._CtxKey, slog.New(l.HandlerOptions.Get()(f, &slog.HandlerOptions{
		AddSource: l.SourceOptions.Get(),
		Level:     l._LevelVar,
	}))), nil
}
