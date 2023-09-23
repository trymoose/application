package flags

import (
	"context"
	"github.com/jessevdk/go-flags"
	"log/slog"
)

const _ParserFlags = flags.PassDoubleDash

type (
	Info struct {
		Name  string
		Short string
		Long  string
		New   func() any
	}
	Commands interface {
		SubCommands() []Info
	}
	Groups interface {
		SubGroups() []Info
	}
)

type (
	Activatable interface {
		Activate(context.Context) error
	}
	ModifyContext interface {
		ModifyContext(context.Context) (context.Context, error)
	}
)

type (
	Parser struct {
		_Name        string
		_ParserFlags flags.Options
		_Commands    []Info
		_Groups      []Info
		_AddHelp     bool
		_AddLogger   bool
		_ExitCodes   *ExitCodes[int]
		_LogLevel    *slog.LevelVar
	}

	ExitCodes[I ~int] struct {
		OK    I
		Help  I
		Error I
	}
)

func NewParser[I ~int](appName string, codes *ExitCodes[I], logLevel *slog.LevelVar) *Parser {
	p := &Parser{
		_Name:        appName,
		_ParserFlags: _ParserFlags,
		_ExitCodes: &ExitCodes[int]{
			OK:    int(codes.OK),
			Help:  int(codes.Help),
			Error: int(codes.Error),
		},
		_LogLevel: logLevel,
	}
	return p
}

func (p *Parser) AddCommand(c Info)        { p._Commands = append(p._Commands, c) }
func (p *Parser) AddGroup(g Info)          { p._Groups = append(p._Groups, g) }
func (p *Parser) AddHelpGroup()            { p._AddHelp = true }
func (p *Parser) AddLoggerGroup()          { p._AddLogger = true }
func (p *Parser) Parse() (*Parsed, error)  { return _Parse(p) }
func (p *Parser) LogLevel() *slog.LevelVar { return p._LogLevel }
