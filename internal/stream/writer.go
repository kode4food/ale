package stream

import (
	"bufio"
	"io"

	"github.com/kode4food/ale/data"
)

type (
	// Writer is used to emit values to a File
	Writer interface {
		Write(data.Value)
	}

	// Closer is used to close a File
	Closer interface {
		Close()
	}

	// OutputFunc is a callback used to marshal values to a Writer
	OutputFunc func(*bufio.Writer, data.Value)

	wrappedWriter struct {
		writer *bufio.Writer
		output OutputFunc
		closer io.Closer
	}
)

const (
	// WriteKey is key used to write to a Writer
	WriteKey = data.Keyword("write")

	// CloseKey is the key used to close a file
	CloseKey = data.Keyword("close")
)

// NewWriter wraps a Go Writer, coupling it with an output function
func NewWriter(w io.Writer, o OutputFunc) *data.Object {
	wrapped := &wrappedWriter{
		writer: bufio.NewWriter(w),
		output: o,
	}

	pairs := data.Pairs{
		data.NewCons(WriteKey, bindWriter(wrapped)),
	}

	if c, ok := w.(io.Closer); ok {
		wrapped.closer = c
		pairs = append(pairs,
			data.NewCons(CloseKey, bindCloser(wrapped)),
		)
	}
	return data.NewObject(pairs...)
}

func (w *wrappedWriter) Write(v data.Value) {
	w.output(w.writer, v)
	_ = w.writer.Flush()
}

func (w *wrappedWriter) Close() {
	_ = w.writer.Flush()
	_ = w.closer.Close()
}

func stringToBytes(s string) []byte {
	return []byte(s)
}

// StrOutput is the standard string-based output function
func StrOutput(w *bufio.Writer, v data.Value) {
	_, _ = w.Write(stringToBytes(v.String()))
}

func bindWriter(w Writer) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		for _, f := range args {
			w.Write(f)
		}
		return data.Null
	})
}

func bindCloser(c Closer) data.Function {
	return data.Applicative(func(...data.Value) data.Value {
		c.Close()
		return data.Null
	}, 0)
}
