package stream

import (
	"io"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

// Writer is used to emit or write values to a stream
type Writer func(ale.Value)

const (
	// WriteKey is key used to write to a Writer
	WriteKey = data.Keyword("write")

	// CloseKey is the key used to close a file
	CloseKey = data.Keyword("close")

	// EmitKey is the key used to emit to a Channel
	EmitKey = data.Keyword("emit")

	// SequenceKey is the key used to retrieve the Sequence from a Channel
	SequenceKey = data.Keyword("seq")
)

func bindWriter(w Writer) data.Procedure {
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
		for _, f := range args {
			w(f)
		}
		return data.Null
	})
}

func bindCloser(c io.Closer) data.Procedure {
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
		_ = c.Close()
		return data.Null
	}, 0)
}
