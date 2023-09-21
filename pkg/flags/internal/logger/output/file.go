package output

import (
	"io"
	"os"
)

type File struct {
	File *string `long:"log-file" description:"Write to file. Combined with 'log-output' if present." env:"APP_LOGGER_FILE"`
}

func (o *File) Writer() (io.WriteCloser, error) {
	if o.File == nil {
		return &NopWriter{}, nil
	}

	f, err := os.OpenFile(*o.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &_DiscardOnClose{W: f}, nil
}
