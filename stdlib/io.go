package stdlib

import (
	"bufio"
	"io"

	"gitlab.com/kode4food/ale/data"
)

type (
	// Reader is used to retrieve values from a File
	Reader interface {
		data.Sequence
	}

	// Writer is used to emit values to a File
	Writer interface {
		data.Value
		Write(data.Value)
	}

	// Closer is used to close a File
	Closer interface {
		Close()
	}

	// OutputFunc is a callback used to marshal values to a Writer
	OutputFunc func(*bufio.Writer, data.Value)

	// InputFunc is a callback used to unmarshal values from a Reader
	InputFunc func(*bufio.Reader) (data.Value, bool)

	wrappedWriter struct {
		writer *bufio.Writer
		output OutputFunc
	}

	wrappedClosingWriter struct {
		closer io.Closer
		*wrappedWriter
	}
)

// NewReader wraps a Go Reader, coupling it with an input function
func NewReader(r io.Reader, i InputFunc) Reader {
	var resolver LazyResolver
	br := bufio.NewReader(r)

	resolver = func() (data.Value, data.Sequence, bool) {
		if v, ok := i(br); ok {
			return v, NewLazySequence(resolver), true
		}
		return data.Nil, data.EmptyList, false
	}

	return NewLazySequence(resolver)
}

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

// LineInput is the standard single line input function
func LineInput(r *bufio.Reader) (data.Value, bool) {
	l, err := r.ReadBytes('\n')
	if err == nil {
		return data.String(l[0 : len(l)-1]), true
	}
	if err == io.EOF && len(l) > 0 {
		return data.String(l), true
	}
	return data.Nil, false
}

// RuneInput is the standard single rune input function
func RuneInput(r *bufio.Reader) (data.Value, bool) {
	if c, _, err := r.ReadRune(); err == nil {
		return data.String(string(c)), true
	}
	return data.Nil, false
}
