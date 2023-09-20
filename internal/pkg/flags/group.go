package flags

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/trymoose/errors"
)

type (
	Group interface {
		_Parsable
	}
	SubGroups interface {
		Groups() []Group
	}
	GroupParsed interface {
		Parsed(ctx context.Context) error
	}
)

func (p *Parser[I]) AddGroup(g Group) error {
	if p._Parsed {
		return ErrParsed
	}
	return p._AddGroup(p._Parser.Command.Group, g)
}

func (p *Parser[I]) _AddGroup(group *flags.Group, g Group) error {
	c, ok := p._Groups[group]
	if !ok {
		return errors.New("group does not exist")
	}

	gr, err := group.AddGroup(g.Short(), g.Long(), g)
	if err != nil {
		return err
	}
	c.Groups = append(c.Groups, g)
	p._Groups[gr] = c

	return p._AddSubGroups(gr, g)
}

func (p *Parser[I]) _AddSubGroups(group *flags.Group, a _Parsable) error {
	if subgroups, ok := a.(SubGroups); ok {
		for _, g := range subgroups.Groups() {
			if err := p._AddGroup(group, g); err != nil {
				return err
			}
		}
	}
	return nil
}
