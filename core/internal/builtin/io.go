package builtin

import (
	"io"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/stream"
)

// MakeReader wraps the go Reader with an input function
func MakeReader(r io.Reader, i stream.InputFunc) stream.Reader {
	return stream.NewReader(r, i)
}

// MakeWriter wraps the go Writer with an output function
func MakeWriter(w io.Writer, o stream.OutputFunc) data.Object {
	wrapped := stream.NewWriter(w, o)

	pairs := []data.Pair{
		data.NewCons(data.TypeKey, stream.WriterType),
		data.NewCons(stream.WriterKey, wrapped),
		data.NewCons(stream.WriteKey, bindWriter(wrapped)),
	}

	if c, ok := w.(stream.Closer); ok {
		pairs = append(pairs, data.NewCons(stream.CloseKey, bindCloser(c)))
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
