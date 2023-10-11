package stream

import (
	"bufio"
	"io"

	"github.com/kode4food/ale/data"
)

// OutputFunc is a callback used to marshal values to a Writer
type OutputFunc func(*bufio.Writer, data.Value)

// NewWriter wraps a Go Writer, coupling it with an output function
func NewWriter(w io.Writer, o OutputFunc) *data.Object {
	buf := bufio.NewWriter(w)
	writer := func(v data.Value) {
		o(buf, v)
		_ = buf.Flush()
	}

	if c, ok := w.(io.Closer); ok {
		return data.NewObject(
			data.NewCons(WriteKey, bindWriter(writer)),
			data.NewCons(CloseKey, bindCloser(c)),
		)
	}

	return data.NewObject(
		data.NewCons(WriteKey, bindWriter(writer)),
	)
}

// StrOutput is the standard string-based output function
func StrOutput(w *bufio.Writer, v data.Value) {
	_, _ = w.Write([]byte(v.String()))
}
