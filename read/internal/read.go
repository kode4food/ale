package internal

import (
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
)

type (
	FromString     func(env.Namespace, data.String) (data.Sequence, error)
	MustFromString func(env.Namespace, data.String) data.Sequence
	MustTokenizer  func(data.String) data.Sequence
)

func MakeTokenizer(matcher lex.Matcher) parse.Tokenizer {
	return func(src data.String) (data.Sequence, error) {
		return lex.Match(src, matcher), nil
	}
}

func MakeFromString(fn parse.Tokenizer) FromString {
	return func(ns env.Namespace, src data.String) (data.Sequence, error) {
		return parse.FromString(ns, fn, src)
	}

}

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
