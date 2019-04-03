package stdlib

import (
	"bufio"
	"io"

	"gitlab.com/kode4food/ale/api"
)

type (
	// Reader is used to retrieve values from a File
	Reader interface {
		api.Sequence
	}

	// Writer is used to emit values to a File
	Writer interface {
		api.Value
		Write(api.Value)
	}

	// Closer is used to close a File
	Closer interface {
		Close()
	}

	// OutputFunc is a callback used to marshal values to a Writer
	OutputFunc func(*bufio.Writer, api.Value)

	// InputFunc is a callback used to unmarshal values from a Reader
	InputFunc func(*bufio.Reader) (api.Value, bool)

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

	resolver = func() (api.Value, api.Sequence, bool) {
		if v, ok := i(br); ok {
			return v, NewLazySequence(resolver), true
		}
		return api.Nil, api.EmptyList, false
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

func (w *wrappedWriter) Write(v api.Value) {
	w.output(w.writer, v)
	w.writer.Flush()
}

func (w *wrappedClosingWriter) Close() {
	w.writer.Flush()
	w.closer.Close()
}

func (w *wrappedWriter) String() string {
	return api.DumpString(w)
}

func (w *wrappedWriter) Type() api.Name {
	return "writer"
}

func stringToBytes(s string) []byte {
	return []byte(s)
}

// StrOutput is the standard string-based output function
func StrOutput(w *bufio.Writer, v api.Value) {
	w.Write(stringToBytes(v.String()))
}

// LineInput is the standard single line input function
func LineInput(r *bufio.Reader) (api.Value, bool) {
	l, err := r.ReadBytes('\n')
	if err == nil {
		return api.String(l[0 : len(l)-1]), true
	}
	if err == io.EOF && len(l) > 0 {
		return api.String(l), true
	}
	return api.Nil, false
}

// RuneInput is the standard single rune input function
func RuneInput(r *bufio.Reader) (api.Value, bool) {
	if c, _, err := r.ReadRune(); err == nil {
		return api.String(string(c)), true
	}
	return api.Nil, false
}
