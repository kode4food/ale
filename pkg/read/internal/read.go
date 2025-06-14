package internal

import (
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

type (
	FromString     func(env.Namespace, data.String) (data.Sequence, error)
	MustFromString func(env.Namespace, data.String) data.Sequence

	MustTokenizer func(data.String) data.Sequence
)

func MakeMustFromString(fn FromString) MustFromString {
	return func(ns env.Namespace, str data.String) data.Sequence {
		seq, err := fn(ns, str)
		if err != nil {
			panic(err)
		}
		return seq
	}
}

func MakeMustTokenizer(fn parse.Tokenizer) MustTokenizer {
	return func(str data.String) data.Sequence {
		seq, err := fn(str)
		if err != nil {
			panic(err)
		}
		return seq
	}
}
