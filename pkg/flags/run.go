package flags

import (
	"context"
	"github.com/trymoose/debug"
	"github.com/trymoose/errors"
	"log/slog"
	"os"
)

func (p *Parsed) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, _ContextKeyLogger, _DefaultLogger())
	ctx = context.WithValue(ctx, _ContextKeyExit, p._Exit(cancel, p._ExitCodes))
	ctx = context.WithValue(ctx, _ContextKeyArgs, p._Args)

	defer func() {
		var finalErr error
		if r := recover(); r != nil {
			// Don't care about the underlying error, just that it is an error.
			//goland:noinspection GoTypeAssertionOnErrors
			if err, ok := r.(error); ok {
				finalErr = errors.Join(finalErr, err)
			} else {
				finalErr = errors.Join(finalErr, errors.New("%s", r))
			}
		}

		if finalErr != nil {
			ContextLogger(ctx).Error("parsing failed", "error", finalErr)
			ContextExit(ctx, p._ExitCodes.Error)
		}
		ContextExit(ctx, p._ExitCodes.OK)
	}()

	return _Run(p._Activated, ctx)
}

func _Run(a *_Activated, ctx context.Context) (err error) {
	for _, a := range append(a.Groups, a.Command) {
		if ctx, err = a.ModifyContext(ctx); err != nil {
			return err
		}
	}

	for _, a := range append(a.Groups, a.Command) {
		if err = a.Activate(ctx); err != nil {
			return err
		}
	}

	if a.Next != nil {
		return _Run(a.Next, ctx)
	}
	return nil
}

func _DefaultLogger() any {
	if debug.Debug {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	}
	return slog.Default()
}

type _Exiter func(context.Context, int)

func (p *Parsed) _Exit(cancel context.CancelFunc, exit *ExitCodes[int]) _Exiter {
	return func(ctx context.Context, i int) {
		cancel()
		if p._Logger != nil {
			if err := p._Logger.Close(); err != nil {
				ContextLogger(ctx).Error("issue closing logger", "error", err)
				i = exit.Error
			}
		}
		os.Exit(i)
	}
}
