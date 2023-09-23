package help

import (
	"bytes"
	"errors"
	"github.com/jessevdk/go-flags"
	"io"
	"os"
	"os/exec"
)

const (
	Name  = "help"
	Short = "Help options."
	Long  = "Display help."
)

var ErrFailedToExit = errors.New("failed to exit")

type Help struct {
	Help      func(string) error `long:"help" choice:"text" choice:"man" description:"Print application help and exit."`
	QuickHelp func() error       `short:"h" description:"Print application help text and exit."`
}

type Exit struct {
	CodeHelp int
	CodeErr  int
	Exit     func(int)
}

func New(exit *Exit, p *flags.Parser) *Help {
	var buf bytes.Buffer
	printBuf := func() {
		_, _ = io.Copy(os.Stderr, &buf)
		exit.Help()
	}

	return &Help{
		Help: func(s string) error {
			switch s {
			case "man":
				if err := _ManPage(p, &buf); err != nil {
					return err
				}
				printBuf()
			case "text":
				if err := _Text(p, &buf); err != nil {
					return err
				}
				printBuf()
			}
			exit.Err()
			return ErrFailedToExit
		},
		QuickHelp: func() error {
			if err := _Text(p, &buf); err != nil {
				return err
			}
			printBuf()
			exit.Err()
			return ErrFailedToExit
		},
	}
}

func _ManPage(p *flags.Parser, out io.Writer) error {
	var buf bytes.Buffer
	p.WriteManPage(&buf)
	cmd := exec.Command("nroff", "-man")
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

func (e *Exit) Help() { e.Exit(e.CodeHelp) }
func (e *Exit) Err()  { e.Exit(e.CodeErr) }
