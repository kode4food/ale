package stream

import (
	"bufio"
	"io"

	"github.com/kode4food/ale/data"
)

type (
	// Writer is used to emit values to a File
	Writer interface {
		data.Value
		Write(data.Value)
	}

	// OutputFunc is a callback used to marshal values to a Writer
	OutputFunc func(*bufio.Writer, data.Value)

	wrappedWriter struct {
		writer *bufio.Writer
		output OutputFunc
	}

	wrappedClosingWriter struct {
		closer io.Closer
		*wrappedWriter
	}
)

// NewWriter wraps a Go Writer, coupling it with an output function
func NewWriter(w io.Writer, o OutputFunc) Writer {
	wrapped := &wrappedWriter{
		writer: bufio.NewWriter(w),
		output: o,
	}
	if c, ok := w.(io.Closer); ok {
		return &wrappedClosingWriter{
			wrappedWriter: wrapped,
			closer:        c,
		}
	}
	return wrapped
}

func (w *wrappedWriter) Write(v data.Value) {
	w.output(w.writer, v)
	w.writer.Flush()
}

func (w *wrappedClosingWriter) Close() {
	w.writer.Flush()
	w.closer.Close()
}

func (w *wrappedWriter) String() string {
	return data.DumpString(w)
}

func (w *wrappedWriter) Type() data.Name {
	return "writer"
}

func stringToBytes(s string) []byte {
	return []byte(s)
}

// StrOutput is the standard string-based output function
func StrOutput(w *bufio.Writer, v data.Value) {
	w.Write(stringToBytes(v.String()))
}
