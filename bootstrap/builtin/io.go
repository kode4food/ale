package builtin

import (
	"io"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/stdlib"
)

const (
	// WriterKey is the key used to wrap a Writer
	WriterKey = api.Keyword("writer")

	// WriteKey is key used to write to a Writer
	WriteKey = api.Keyword("write")

	// CloseKey is the key used to close a file
	CloseKey = api.Keyword("close")
)

const typeKey = api.Keyword("type")

var writerPrototype = api.Object{
	typeKey: api.Name("writer"),
}

// MakeReader wraps the go Reader with an input function
func MakeReader(r io.Reader, i stdlib.InputFunc) stdlib.Reader {
	return stdlib.NewReader(r, i)
}

// MakeWriter wraps the go Writer with an output function
func MakeWriter(w io.Writer, o stdlib.OutputFunc) api.Object {
	wrapped := stdlib.NewWriter(w, o)

	wrapper := api.Object{
		WriterKey: wrapped,
		WriteKey:  bindWriter(wrapped),
	}

	if c, ok := w.(stdlib.Closer); ok {
		wrapper[CloseKey] = bindCloser(c)
	}

	return writerPrototype.Extend(wrapper)
}

func bindWriter(w stdlib.Writer) *api.Function {
	return api.ApplicativeFunction(func(args ...api.Value) api.Value {
		for _, f := range args {
			w.Write(f)
		}
		return api.Nil
	})
}

func bindCloser(c stdlib.Closer) *api.Function {
	return api.ApplicativeFunction(func(args ...api.Value) api.Value {
		c.Close()
		return api.Nil
	})
}
