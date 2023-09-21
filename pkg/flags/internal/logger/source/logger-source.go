//go:build !debug

package source

type SourceOptions struct {
	Source bool `long:"log-source" description:"Logs the source file and line number." env:"APP_LOGGER_SOURCE"`
}

func (o SourceOptions) Get() bool { return o.Source }
