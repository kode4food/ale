package builtin

import (
	"io"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/stdlib"
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
func MakeReader(r io.Reader, i stdlib.InputFunc) stdlib.Reader {
	return stdlib.NewReader(r, i)
}

// MakeWriter wraps the go Writer with an output function
func MakeWriter(w io.Writer, o stdlib.OutputFunc) data.Object {
	wrapped := stdlib.NewWriter(w, o)

	res := data.Object{
		data.TypeKey: WriterType,
		WriterKey:    wrapped,
		WriteKey:     bindWriter(wrapped),
	}

	if c, ok := w.(stdlib.Closer); ok {
		res[CloseKey] = bindCloser(c)
	}

	return res
}

func bindWriter(w stdlib.Writer) data.Call {
	return func(args ...data.Value) data.Value {
		for _, f := range args {
			w.Write(f)
		}
		return data.Nil
	}
}

func bindCloser(c stdlib.Closer) data.Call {
	return func(args ...data.Value) data.Value {
		c.Close()
		return data.Nil
	}
}
