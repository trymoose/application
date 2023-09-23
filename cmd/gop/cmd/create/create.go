package create

import (
	"context"
	_ "embed"
	app "github.com/trymoose/application"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed default.gotmpl
var DefaultTemplate []byte

func init() {
	app.AddCommand(app.Info{
		Name:  "create",
		Short: "Create a new go project.",
		Long:  "Create a new go project",
		New:   func() any { return new(Options) },
	})
}

type Options struct {
	SkipInitialCommit bool `short:"s" long:"skip-initial-commit" description:"Do not initialize the current directory as a git repo."`
	Positional        struct {
		Name string `positional-arg-name:"name"`
	} `positional-args:"true" required:"true"`
}

func (o *Options) Activate(ctx context.Context) error {
	if err := Run(ctx, "go", "mod", "init", o.Positional.Name); err != nil {
		return err
	}

	if err := os.Mkdir("cmd", os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join("cmd", "main.go"), DefaultTemplate, os.ModePerm); err != nil {
		return err
	}

	if err := Run(ctx, "go", "get", "."); err != nil {
		return err
	}

	if o.SkipInitialCommit {
		return nil
	}

	if err := Run(ctx, "git", "init"); err != nil {
		return err
	}

	if err := Run(ctx, "git", "add", "."); err != nil {
		return err
	}

	return Run(ctx, "git", "commit", "-m", "initial commit")
}

func Run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
