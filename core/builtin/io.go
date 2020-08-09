package builtin

import (
	"io"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/stream"
)

const (
	// WriterType is the type name for a writer
	WriterType = data.String("writer")

	// WriterKey is the key used to wrap a Writer
	WriterKey = data.Keyword("writer")

	// WriteKey is key used to write to a Writer
	WriteKey = data.Keyword("write")

	// CloseKey is the key used to close a file
	CloseKey = data.Keyword("close")
)

// MakeReader wraps the go Reader with an input function
func MakeReader(r io.Reader, i stream.InputFunc) stream.Reader {
	return stream.NewReader(r, i)
}

// MakeWriter wraps the go Writer with an output function
func MakeWriter(w io.Writer, o stream.OutputFunc) data.Object {
	wrapped := stream.NewWriter(w, o)

	res := data.Object{
		data.TypeKey: WriterType,
		WriterKey:    wrapped,
		WriteKey:     bindWriter(wrapped),
	}

	if c, ok := w.(stream.Closer); ok {
		res[CloseKey] = bindCloser(c)
	}

	return res
}

func bindWriter(w stream.Writer) data.Call {
	return func(args ...data.Value) data.Value {
		for _, f := range args {
			w.Write(f)
		}
		return data.Nil
	}
}

func bindCloser(c stream.Closer) data.Call {
	return func(args ...data.Value) data.Value {
		c.Close()
		return data.Nil
	}
}
