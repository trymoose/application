package help

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"os/exec"
)

var ErrFailedToExit = errors.New("failed to exit")

type Help struct {
	Help      func(string) error `long:"help" choice:"text" choice:"man" description:"Print application help and exit."`
	QuickHelp func() error       `short:"h" description:"Print application help text and exit."`
}

func (*Help) Name() string  { return "help" }
func (*Help) Short() string { return "Help options." }
func (*Help) Long() string  { return "Display help." }

type Exit[I ~int] struct {
	CodeHelp I
	CodeErr  I
	Exit     func(context.Context, I)
}

func New[I ~int](ctx context.Context, exit *Exit[I], p *flags.Parser) *Help {
	var buf bytes.Buffer
	printBuf := func() {
		fmt.Println(buf.String())
		exit.Help(ctx)
	}

	return &Help{
		Help: func(s string) error {
			switch s {
			case "man":
				if err := _ManPage(ctx, p, &buf); err != nil {
					return err
				}
				printBuf()
			case "text":
				if err := _Text(p, &buf); err != nil {
					return err
				}
				printBuf()
			}
			exit.Err(ctx)
			return ErrFailedToExit
		},
		QuickHelp: func() error {
			if err := _Text(p, &buf); err != nil {
				return err
			}
			printBuf()
			exit.Err(ctx)
			return ErrFailedToExit
		},
	}
}

func _ManPage(ctx context.Context, p *flags.Parser, out io.Writer) error {
	var buf bytes.Buffer
	p.WriteManPage(&buf)
	cmd := exec.CommandContext(ctx, "nroff", "-man")
	cmd.Stdin = &buf
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func _Text(p *flags.Parser, out io.Writer) error {
	p.WriteHelp(out)
	return nil
}

func (e *Exit[I]) Help(ctx context.Context) { e.Exit(ctx, e.CodeHelp) }
func (e *Exit[I]) Err(ctx context.Context)  { e.Exit(ctx, e.CodeErr) }
