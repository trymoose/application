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

func addCommands(ctx context.Context, cmd *flags.Parser, subcommands []Command, curr *cmdr) {
	for _, c := range subcommands {
		name := c.Name()
		l, s := c.Description()
		errors.Get(cmd.AddCommand(name, l, s, &cmdr{
			ctx:     ctx,
			cmd:     c,
			prevCmd: curr,
		}))
	}
}

func AppName() string { return _AppName }
