package output

import "io"

type OutputOptions struct {
	File
	Stdio
}

func (o *OutputOptions) Writer() (io.WriteCloser, error) {
	f, err := o.File.Writer()
	if err != nil {
		return nil, err
	}

	return MultiWriteCloser(f, &_WriteClosers{W: o.Stdio.Writer()}), nil
}
