package application

import (
	"context"
)

type cmdr struct {
	ctx     context.Context
	exec    func(context.Context) error
	prevCmd *cmdr
}

func (c *cmdr) Execute([]string) error {
	if c.prevCmd != nil {
		if err := c.prevCmd.Execute(nil); err != nil {
			return err
		}
	}
	return c.exec(c.ctx)
}

var _NOOPExec = func(context.Context) error { return nil }
