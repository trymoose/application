package application

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/trymoose/errors"
)

type cmdr struct {
	ctx     context.Context
	cmd     Command
	prevCmd *cmdr
}

func (c *cmdr) Execute(args []string) error {
	// Create flag parser for this command
	cmd := flags.NewParser(c.cmd, flags.Default)
	// No op command function
	exec := func(context.Context) error { return nil }
	// This command has subcommands, add them
	if sc, ok := c.cmd.(Subcommander); ok {
		addCommands(c.ctx, cmd, sc.Subcommands(), c)
	} else {
		exec = c.cmd.Execute
	}

	// This command has subcommands, but would like to be run too
	if _, ok := c.cmd.(ExecSubcommander); ok {
		exec = c.cmd.Execute
	}
	errors.Get(cmd.ParseArgs(args))
	return exec(c.ctx)
}
