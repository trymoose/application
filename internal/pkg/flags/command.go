package flags

import (
	"context"
	"github.com/jessevdk/go-flags"
)

type (
	Command interface {
		_Parsable
	}
	RunCommand interface {
		Run(context.Context) error
	}
	SubCommands interface {
		SubCommand() []Command
	}
)

func (p *Parser[I]) AddCommand(c Command) error {
	if p._Parsed {
		return ErrParsed
	}
	return p._AddCommand(p._Parser.Command, c)
}

func (p *Parser[I]) _AddCommand(cmd *flags.Command, c Command) error {
	cmd, err := cmd.AddCommand(c.Name(), c.Short(), c.Long(), c)
	if err != nil {
		return err
	}

	if err := p._AddSubGroups(cmd.Group, c); err != nil {
		return err
	}

	p._Commands[cmd] = &_Command{Command: c}
	if sub, ok := c.(SubCommands); ok {
		for _, sub := range sub.SubCommand() {
			if err := p._AddCommand(cmd, sub); err != nil {
				return err
			}
		}
	}
	return nil
}
