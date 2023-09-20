package app

import (
	"context"
	"github.com/trymoose/application/internal/pkg/flags"
	"github.com/trymoose/debug"
	"github.com/trymoose/errors"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
)

const Debug = debug.Debug

type (
	Command     = flags.Command
	RunCommand  = flags.RunCommand
	SubCommands = flags.SubCommands
)

type (
	Group       = flags.Group
	SubGroups   = flags.SubGroups
	GroupParsed = flags.GroupParsed
)

type ModCtx = flags.ModCtx

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
	_Parser.Parse(ctx)
}

// AddGroup adds a group to the root
func AddGroup(g Group) { errors.Check(_Parser.AddGroup(g)) }

// AddCommand adds a command to the application.
// Commands must implement [Command]
// Optionally commands can implement:
// - [RunCommand]
// - [SubCommands]
// - [ModCtx]
func AddCommand(c Command) { errors.Check(_Parser.AddCommand(c)) }

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
