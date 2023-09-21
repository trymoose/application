package output

import (
	"github.com/trymoose/application/internal/pkg/misc"
	"github.com/trymoose/errors"
	"io"
	"sync"
)

type _DiscardOnClose struct {
	W       io.Writer
	_Closed bool
	_Lock   sync.RWMutex
}

func (w *_DiscardOnClose) Write(b []byte) (int, error) {
	w._Lock.RLock()
	defer w._Lock.RUnlock()
	if w._Closed {
		return len(b), nil
	}
	return w.W.Write(b)
}

func (w *_DiscardOnClose) Close() error {
	w._Lock.Lock()
	defer w._Lock.Unlock()
	if w._Closed {
		return nil
	}

	if cl, ok := w.W.(io.Closer); ok {
		return cl.Close()
	}
	return nil
}

type _NoopClose struct{ W io.Writer }

func (c *_NoopClose) Write(p []byte) (n int, err error) { return c.W.Write(p) }
func (c *_NoopClose) Close() error                      { return nil }

type _WriteClosers struct {
	W  io.Writer
	Cl func() error
}

func (wc *_WriteClosers) Write(p []byte) (n int, err error) { return wc.W.Write(p) }
func (wc *_WriteClosers) Close() error {
	if wc.Cl == nil {
		return nil
	}
	return wc.Cl()
}

func MultiWriteCloser(w ...io.WriteCloser) io.WriteCloser {
	wc := &_WriteClosers{W: io.MultiWriter(misc.MapSlice(w, func(w io.WriteCloser) io.Writer { return w })...)}

	wc.Cl = func() error {
		wc.Cl = func() error { return nil }
		return nil
	}

	for _, w := range w {
		fn := wc.Cl
		wn := w
		wc.Cl = func() error {
			return errors.Join(fn(), wn.Close())
		}
	}

	return wc
}

type NopWriter struct{ io.Writer }

func (nw NopWriter) Write(b []byte) (int, error) {
	if nw.Writer == nil {
		return len(b), nil
	}
	return nw.Write(b)
}

func (nw NopWriter) Close() error {
	if w, ok := nw.Writer.(io.WriteCloser); ok {
		return w.Close()
	}
	return nil
}
