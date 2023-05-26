package application

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/trymoose/debug"
	"github.com/trymoose/errors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const Debug = debug.Debug

func init() {
	if Debug {
		go _PrintDebug()
	}
}

var _RegisteredSubcommands []Command

func RegisterSubcommand(cmd Command) {
	_RegisteredSubcommands = append(_RegisteredSubcommands, cmd)
}

type Command interface {
	Name() string
	Description() (long, short string)
	Execute(ctx context.Context) error
}

type Subcommander interface {
	Subcommands() []Command
}

type ExecSubcommander interface {
	Executable()
}

var (
	wg     *errgroup.Group
	ctx    context.Context
	cancel context.CancelFunc
)

var _AppName string

func init() {
	ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	_SetExitCancel(cancel)
	wg, ctx = errgroup.WithContext(ctx)

	_AppName = os.Args[0]
	_, _AppName = filepath.Split(_AppName)
	_AppName = strings.TrimSuffix(_AppName, filepath.Ext(_AppName))
}

func Main() {
	defer cancel()
	cmd := flags.NewNamedParser(AppName(), flags.Default)
	addCommands(ctx, cmd, _RegisteredSubcommands, nil)
	_, err := cmd.ParseArgs(os.Args[1:])
	if err, ok := errors.To[*flags.Error](err); ok && err.Type == flags.ErrHelp {
		return
	}
	errors.Check(wg.Wait())
}

type Parser interface {
	AddCommand(command string, shortDescription string, longDescription string, data interface{}) (*flags.Command, error)
	AddGroup(shortDescription string, longDescription string, data interface{}) (*flags.Group, error)
}

func addCommands(ctx context.Context, cmd Parser, subcommands []Command, curr *cmdr) {
	for _, c := range subcommands {
		curr := &cmdr{
			ctx:     ctx,
			exec:    c.Execute,
			prevCmd: curr,
		}

		sub := errors.Get(cmd.AddCommand(c.Name(), "", "", curr))
		errors.Get(sub.AddGroup(func(l, s string) (string, string, any) { return s, l, c }(c.Description())))

		if sc, ok := c.(Subcommander); ok {
			if _, ok := c.(ExecSubcommander); !ok {
				curr.exec = _NOOPExec
			}
			addCommands(ctx, sub, sc.Subcommands(), curr)
		}
	}
}

func AppName() string { return _AppName }
