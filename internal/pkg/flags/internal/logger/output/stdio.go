package output

import (
	"github.com/trymoose/application/internal/pkg/misc"
	"io"
	"os"
)

type Stdio struct {
	Stdio []Output `long:"log-output" choice:"stdout" choice:"stderr" description:"Output to [os.Stdout] and/or [os.Stderr]." default:"stderr" env:"APP_LOGGER_STDIO"`
}

type Output string

const (
	Stdout Output = "stdout"
	Stderr Output = "stderr"
)

func (o Output) Write(b []byte) (int, error) {
	switch o {
	case Stdout:
		return os.Stdout.Write(b)
	case Stderr:
		return os.Stderr.Write(b)
	default:
		return io.Discard.Write(b)
	}
}

func (s *Stdio) Writer() io.Writer {
	return io.MultiWriter(misc.MapSlice(s.Stdio, func(e Output) io.Writer { return e })...)
}
