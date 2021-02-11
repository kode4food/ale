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

	pairs := []data.Pair{
		data.NewCons(data.TypeKey, WriterType),
		data.NewCons(WriterKey, wrapped),
		data.NewCons(WriteKey, bindWriter(wrapped)),
	}

	if c, ok := w.(stream.Closer); ok {
		pairs = append(pairs, data.NewCons(CloseKey, bindCloser(c)))
	}

	return data.NewObject(pairs...)
}

func bindWriter(w stream.Writer) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		for _, f := range args {
			w.Write(f)
		}
		return data.Nil
	})
}

func bindCloser(c stream.Closer) data.Function {
	return data.Applicative(func(_ ...data.Value) data.Value {
		c.Close()
		return data.Nil
	}, 0)
}
