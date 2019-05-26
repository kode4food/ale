package builtin

import (
	"io"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/stdlib"
)

const (
	// WriterKey is the key used to wrap a Writer
	WriterKey = data.Keyword("writer")

	// WriteKey is key used to write to a Writer
	WriteKey = data.Keyword("write")

	// CloseKey is the key used to close a file
	CloseKey = data.Keyword("close")
)

const typeKey = data.Keyword("type")

var writerPrototype = data.Object{
	typeKey: data.Name("writer"),
}

// MakeReader wraps the go Reader with an input function
func MakeReader(r io.Reader, i stdlib.InputFunc) stdlib.Reader {
	return stdlib.NewReader(r, i)
}

// MakeWriter wraps the go Writer with an output function
func MakeWriter(w io.Writer, o stdlib.OutputFunc) data.Object {
	wrapped := stdlib.NewWriter(w, o)

	wrapper := data.Object{
		WriterKey: wrapped,
		WriteKey:  bindWriter(wrapped),
	}

	if c, ok := w.(stdlib.Closer); ok {
		wrapper[CloseKey] = bindCloser(c)
	}

	return writerPrototype.Extend(wrapper)
}

func bindWriter(w stdlib.Writer) *data.Function {
	return data.ApplicativeFunction(func(args ...data.Value) data.Value {
		for _, f := range args {
			w.Write(f)
		}
		return data.Nil
	})
}

func bindCloser(c stdlib.Closer) *data.Function {
	return data.ApplicativeFunction(func(args ...data.Value) data.Value {
		c.Close()
		return data.Nil
	})
}
