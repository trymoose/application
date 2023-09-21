package app

import (
	"context"
	"fmt"
	"github.com/trymoose/application/pkg/flags"
	"github.com/trymoose/debug"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
)

const Debug = debug.Debug

type (
	Info          = flags.Info
	Commands      = flags.Commands
	Groups        = flags.Groups
	Activatable   = flags.Activatable
	ModifyContext = flags.ModifyContext
)

var _Parser = flags.NewParser(Name(), &flags.ExitCodes[ExitCode]{
	OK:    ExitCodeOK,
	Error: ExitCodeError,
	Help:  ExitCodeHelp,
})

func AddHelpGroup()   { _Parser.AddHelpGroup() }
func AddLoggerGroup() { _Parser.AddLoggerGroup() }
func AddDefaultGroups() {
	AddHelpGroup()
	AddLoggerGroup()
}

func Name() string                             { return filepath.Base(os.Args[0]) }
func Args(ctx context.Context) []string        { return flags.ContextArgs(ctx) }
func Exit[T ~int](ctx context.Context, code T) { flags.ContextExit(ctx, code) }
func Logger(ctx context.Context) *slog.Logger  { return flags.ContextLogger(ctx) }
func LoggerWith(ctx context.Context, args ...any) context.Context {
	return flags.ContextLoggerWith(ctx, args...)
}

func Main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()
	if p, err := _Parser.Parse(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse application args: %v\n", err)
		os.Exit(int(ExitCodeError))
	} else if err := p.Run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to run application: %v\n", err)
		os.Exit(int(ExitCodeError))
	}
}

// AddGroup adds a group to the root
func AddGroup(g Info) { _Parser.AddGroup(g) }

// AddCommand adds a command to the application.
// Commands must implement [Command]
// Optionally commands can implement:
// - [RunCommand]
// - [SubCommands]
// - [ModCtx]
func AddCommand(c Info) { _Parser.AddCommand(c) }

// ExitCode is a code passed on exit to the os.
type ExitCode int

const (
	// ExitCodeOK is the zero exit code
	ExitCodeOK ExitCode = iota
	// ExitCodeError means some error happened in your application
	ExitCodeError
	// ExitCodeHelp is used when the application prints its help.
	ExitCodeHelp
)
