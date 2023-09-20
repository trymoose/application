package flags

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/trymoose/application/internal/pkg/flags/internal/help"
	"github.com/trymoose/application/internal/pkg/flags/internal/logger"
	"github.com/trymoose/debug"
	"github.com/trymoose/errors"
	"log/slog"
	"os"
)

const _ParserFlags = flags.PassDoubleDash

type _Parsable interface {
	Name() string
	Short() string
	Long() string
}

var ErrParsed = errors.New("flags already parsed")

type (
	Parser[I ~int] struct {
		_Parser    *flags.Parser
		_Commands  map[*flags.Command]*_Command
		_Groups    map[*flags.Group]*_Command
		_Parsed    bool
		_AddHelp   bool
		_AddLogger bool
		_ExitCodes *ExitCodes[I]
	}

	_Command struct {
		Command Command
		Groups  []Group
	}
)

func NewParser[I ~int](appName string, codes *ExitCodes[I]) *Parser[I] {
	p := &Parser[I]{
		_Parser:    flags.NewNamedParser(appName, _ParserFlags),
		_Commands:  map[*flags.Command]*_Command{},
		_Groups:    map[*flags.Group]*_Command{},
		_ExitCodes: codes,
	}
	// Use command for proper type check in [Parser._Run]
	p._Commands[p._Parser.Command] = &_Command{Command: _ParserCommand{}}
	p._Groups[p._Parser.Command.Group] = p._Commands[p._Parser.Command]

	p._Parser.CommandHandler = func(cmd flags.Commander, args []string) error { return nil }
	return p
}

type _Exiter func(context.Context, int)

func (p *Parser[I]) Parse(ctx context.Context) (finalErr error) {
	ctx = context.WithValue(ctx, _ContextKeyLogger, _DefaultLogger())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	closeLogger := func() error { return nil }
	if p._AddLogger && !p._Parsed {
		lg := logger.New(_ContextKeyLogger)
		if err := p.AddGroup(lg); err != nil {
			return err
		}
		closeLogger = lg.Close
	}

	ctx = context.WithValue(ctx, _ContextKeyExit, _Exiter(func(ctx context.Context, i int) {
		cancel()
		if err := closeLogger(); err != nil {
			ContextLogger(ctx).Error("issue closing logger", "error", err)
			i = int(p._ExitCodes.Error)
		}
		os.Exit(i)
	}))

	defer func() {
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

	if p._AddHelp && !p._Parsed {
		if err := p.AddGroup(help.New[I](ctx, &help.Exit[I]{
			CodeHelp: p._ExitCodes.Help,
			CodeErr:  p._ExitCodes.Error,
			Exit:     ContextExit[I],
		}, p._Parser)); err != nil {
			return err
		}
	}

	if p._Parsed {
		return ErrParsed
	}
	p._Parsed = true

	args, err := p._Parser.Parse()
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, _ContextKeyArgs, args)
	return p._Run(ctx, p._Parser.Command, &ctx)
}

func (p *Parser[I]) _Run(ctx context.Context, cmd *flags.Command, finalCtx *context.Context) (err error) {
	*finalCtx = ctx
	if cmd == nil {
		return nil
	}

	command, ok := p._Commands[cmd]
	if !ok {
		return errors.New("command(%q) not registered", cmd.Name)
	}

	for _, group := range command.Groups {
		if g, ok := group.(ModCtx); ok {
			if ctx, err = g.ModCtx(ctx); err != nil {
				return err
			}
		}

		if g, ok := group.(GroupParsed); ok {
			if err := g.Parsed(ctx); err != nil {
				return err
			}
		}
	}

	if cmd, ok := command.Command.(ModCtx); ok {
		if ctx, err = cmd.ModCtx(ctx); err != nil {
			return err
		}
	}

	if cmd, ok := command.Command.(RunCommand); ok {
		if err := cmd.Run(ctx); err != nil {
			return err
		}
	}
	return p._Run(ctx, cmd.Active, finalCtx)
}

func (p *Parser[I]) Raw() *flags.Parser { return p._Parser }

type _ParserCommand struct{}

func (_ParserCommand) Name() string  { return "" }
func (_ParserCommand) Short() string { return "" }
func (_ParserCommand) Long() string  { return "" }

func (p *Parser[I]) AddHelpGroup()   { p._AddHelp = true }
func (p *Parser[I]) AddLoggerGroup() { p._AddLogger = true }

type ExitCodes[I ~int] struct {
	OK    I
	Help  I
	Error I
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
