package flags

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/trymoose/application/pkg/flags/internal/help"
	"github.com/trymoose/application/pkg/flags/internal/logger"
	"github.com/trymoose/errors"
	"log/slog"
	"os"
)

type Parsed struct {
	_Parser    *flags.Parser
	_Args      []string
	_Logger    *logger.Logger
	_Activated *_Activated
	_ExitCodes *ExitCodes[int]
	_LogLevel  *slog.LevelVar
}

type (
	_FlgCmd           = *flags.Command
	_FlgGrp           = *flags.Group
	_Flg              interface{ _FlgCmd | _FlgGrp }
	_ParseMap[K _Flg] map[K]*_Interfaces
	_Interfaces       struct {
		Name string
		_Act Activatable
		_Mod ModifyContext
	}
	_Activated struct {
		Command *_Interfaces
		Groups  []*_Interfaces
		Next    *_Activated
	}
)

func _Parse(info *Parser) (_ *Parsed, finalErr error) {
	p := Parsed{
		_Parser:    flags.NewNamedParser(info._Name, info._ParserFlags),
		_ExitCodes: info._ExitCodes,
		_LogLevel:  info._LogLevel,
	}
	p._Parser.CommandHandler = func(flags.Commander, []string) error { return nil }

	commands := _ParseMap[_FlgCmd]{}
	groups := _ParseMap[_FlgGrp]{}
	_AddRoot(info._Name, p._Parser, commands, groups)

	if err := errors.Join(
		_AddGroups(info._Name, p._AddBuiltInGroups(info._Groups, info._AddLogger, info._AddHelp), p._Parser.Group, groups),
		_AddCommands(info._Name, info._Commands, p._Parser.Command, commands, groups),
	); err != nil {
		return nil, err
	}

	args, err := p._Parser.Parse()
	if err != nil {
		flagErr, ok := errors.To[*flags.Error](err)
		if !ok {
			return nil, errors.New("non flag error returned: %w", err)
		} else if flagErr.Type == flags.ErrCommandRequired {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return nil, errors.Join(p._CreateHelp().QuickHelp(), err)
		}
		return nil, err
	}
	p._Args = args
	p._Activated = p._CollectActivated(p._Parser.Command, commands, groups)
	return &p, nil
}

func (p *Parsed) _CollectActivated(cmd *flags.Command, commands _ParseMap[_FlgCmd], groups _ParseMap[_FlgGrp]) *_Activated {
	var a _Activated
	if c, ok := commands[cmd]; ok {
		a.Command = c
	}

	for _, g := range cmd.Groups() {
		if g, ok := groups[g]; ok {
			a.Groups = append(a.Groups, g)
		}
	}

	if cmd.Active != nil {
		a.Next = p._CollectActivated(cmd.Active, commands, groups)
	}
	return &a
}

func (p *Parsed) _AddBuiltInGroups(groups []Info, addLogger, addHelp bool) []Info {
	if addLogger {
		p._Logger = logger.New(_ContextKeyLogger, p._LogLevel)
		groups = append(groups, p._LoggerGroup())
	}

	if addHelp {
		groups = append(groups, p._HelpGroup())
	}
	return groups
}

func (p *Parsed) _HelpGroup() Info {
	return Info{
		Name:  help.Name,
		Short: help.Short,
		Long:  help.Long,
		New:   func() any { return p._CreateHelp() },
	}
}

func (p *Parsed) _CreateHelp() *help.Help {
	return help.New(&help.Exit{
		CodeHelp: p._ExitCodes.Help,
		CodeErr:  p._ExitCodes.Error,
		Exit:     os.Exit,
	}, p._Parser)
}

func (p *Parsed) _LoggerGroup() Info {
	return Info{
		Name:  logger.Name,
		Short: logger.Short,
		Long:  logger.Long,
		New:   func() any { return p._Logger },
	}
}

func _AddRoot(name string, p *flags.Parser, cmds _ParseMap[_FlgCmd], grps _ParseMap[_FlgGrp]) {
	cmd := &_Interfaces{Name: name}
	cmds[p.Command] = cmd
	grps[p.Group] = cmd
}

func _AddCommand(command *flags.Command, info *Info, commands _ParseMap[_FlgCmd], groups _ParseMap[_FlgGrp]) error {
	v := info.New()
	command, err := command.AddCommand(info.Name, info.Short, info.Long, v)
	if err != nil {
		return errors.New("failed to add command %q: %w", info.Name, err)
	}

	commands[command] = _NewInterfaces(info.Name, v)
	return errors.Join(
		_AddGroups(info.Name, _GetSubCommands(v), command.Group, groups),
		_AddCommands(info.Name, _GetSubGroups(v), command, commands, groups),
	)
}

func _GetSubGroups(a any) []Info {
	if a, ok := a.(Groups); ok {
		return a.SubGroups()
	}
	return nil
}

func _GetSubCommands(a any) []Info {
	if a, ok := a.(Commands); ok {
		return a.SubCommands()
	}
	return nil
}

func _AddCommands(name string, commands []Info, command *flags.Command, cmds _ParseMap[_FlgCmd], grps _ParseMap[_FlgGrp]) error {
	for _, cmd := range commands {
		if err := _AddCommand(command, &cmd, cmds, grps); err != nil {
			return errors.New("failed to add command %q to command %q: %w", cmd.Name, name, err)
		}
	}
	return nil
}

func _AddGroup(group *flags.Group, info *Info, grps _ParseMap[_FlgGrp]) error {
	v := info.New()
	group, err := group.AddGroup(info.Short, info.Long, v)
	if err != nil {
		return errors.New("failed to add group %q: %w", info.Name, err)
	}

	grps[group] = _NewInterfaces(info.Name, v)
	return _AddGroups(info.Name, _GetSubGroups(v), group, grps)
}

func _AddGroups(name string, groups []Info, group *flags.Group, grps _ParseMap[_FlgGrp]) error {
	for _, grp := range groups {
		if err := _AddGroup(group, &grp, grps); err != nil {
			return errors.New("failed to add group %q to group %q: %w", grp.Name, name, err)
		}
	}
	return nil
}

func _NewInterfaces(name string, v any) *_Interfaces {
	var in _Interfaces
	in.Name = name
	in._Mod, _ = v.(ModifyContext)
	in._Act, _ = v.(Activatable)
	return &in
}

func (in *_Interfaces) ModifyContext(ctx context.Context) (context.Context, error) {
	if in != nil && in._Mod != nil {
		if c, err := in._Mod.ModifyContext(ctx); err != nil {
			return nil, errors.New("(%s).ModifyContext(%T) failed: %w", in.Name, ctx, err)
		} else {
			return c, nil
		}
	}
	return ctx, nil
}

func (in *_Interfaces) Activate(ctx context.Context) error {
	if in != nil && in._Act != nil {
		if err := in._Act.Activate(ctx); err != nil {
			return errors.New("(%s).Activate(%T) failed: %w", in.Name, ctx, err)
		}
	}
	return nil
}
