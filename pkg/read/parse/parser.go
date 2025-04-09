package parse

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/read/lex"
)

// parser is a stateful iteration interface for a Token stream that is piloted
// by the FromLexer function and exposed as a LazySequence
type parser struct {
	seq   data.Sequence
	token *lex.Token
}

const (
	// ErrPrefixedNotPaired is raised when the parser encounters the end of the
	// stream without being able to complete a paired element, such as a quote
	ErrPrefixedNotPaired = "end of file reached before completing %s"

	// ErrUnexpectedDot is raised when the parser encounters a dot in the
	// stream when it isn't part of an open list
	ErrUnexpectedDot = "encountered '" + lang.Dot + "' with no open list"

	// ErrInvalidListSyntax is raised when the parse encounters a misplaced dot
	// when parsing an open list
	ErrInvalidListSyntax = "invalid list syntax"

	// ErrListNotClosed is raised when the parser encounters the end of the
	// stream while an open list is still being parsed
	ErrListNotClosed = "end of file reached with open list"

	// ErrUnmatchedListEnd is raised when a list-end character is encountered
	// in the stream when no open list is being parsed
	ErrUnmatchedListEnd = "encountered '" + lang.ListEnd +
		"' with no open list"

	// ErrVectorNotClosed is raised when the parser encounters the end of the
	// stream while an open vector is still being parsed
	ErrVectorNotClosed = "end of file reached with open vector"

	// ErrUnmatchedVectorEnd is raised when a vector-end character is
	// encountered in the stream when no open vector is being parsed
	ErrUnmatchedVectorEnd = "encountered '" + lang.VectorEnd +
		"' with no open vector"

	// ErrObjectNotClosed is raised when the parser encounters the end of the
	// stream while an open object is still being parsed
	ErrObjectNotClosed = "end of file reached with open object"

	// ErrUnmatchedObjectEnd is raised when an object-end character is
	// encountered in the stream when no open object is being parsed
	ErrUnmatchedObjectEnd = "encountered '" + lang.ObjectEnd +
		"' with no open object"
)

var (
	quoteSym    = env.RootSymbol("quote")
	syntaxSym   = env.RootSymbol("syntax-quote")
	unquoteSym  = env.RootSymbol("unquote")
	splicingSym = env.RootSymbol("unquote-splicing")

	specialNames = map[data.String]data.Value{
		data.TrueLiteral:  data.True,
		data.FalseLiteral: data.False,
	}

	collectionErrors = map[lex.TokenType]string{
		lex.VectorEnd: ErrVectorNotClosed,
		lex.ObjectEnd: ErrObjectNotClosed,
	}
)

func newParser(lexer data.Sequence) *parser {
	return &parser{
		seq: lexer,
	}
}

func (r *parser) nextValue() (data.Value, bool, error) {
	t, err := r.nextToken()
	if err != nil {
		return nil, false, r.maybeWrap(err)
	}
	if t != nil {
		v, err := r.value(t)
		if err != nil {
			return nil, false, r.maybeWrap(err)

		}
		return v, true, nil
	}
	return nil, false, nil
}

func (r *parser) nextToken() (*lex.Token, error) {
	token, seq, ok := r.seq.Split()
	if !ok {
		return nil, nil
	}
	r.token = token.(*lex.Token)
	r.seq = seq
	if r.token.Type() == lex.Error {
		return nil, r.error(data.ToString(r.token.Value()))
	}
	return r.token, nil
}

func (r *parser) maybeWrap(err error) error {
	if t := r.token; t != nil {
		return t.WrapError(err)
	}
	return err
}

func (r *parser) value(t *lex.Token) (data.Value, error) {
	switch t.Type() {
	case lex.QuoteMarker:
		return r.prefixed(quoteSym)
	case lex.SyntaxMarker:
		return r.prefixed(syntaxSym)
	case lex.UnquoteMarker:
		return r.prefixed(unquoteSym)
	case lex.SpliceMarker:
		return r.prefixed(splicingSym)
	case lex.ListStart:
		return r.list()
	case lex.VectorStart:
		return r.vector()
	case lex.ObjectStart:
		return r.object()
	case lex.Keyword:
		return r.keyword(), nil
	case lex.Identifier:
		return r.identifier()
	case lex.ListEnd:
		return nil, r.error(ErrUnmatchedListEnd)
	case lex.VectorEnd:
		return nil, r.error(ErrUnmatchedVectorEnd)
	case lex.ObjectEnd:
		return nil, r.error(ErrUnmatchedObjectEnd)
	case lex.Dot:
		return nil, r.error(ErrUnexpectedDot)
	default:
		return t.Value(), nil
	}
}

func (r *parser) prefixed(s data.Symbol) (data.Value, error) {
	v, ok, err := r.nextValue()
	if err != nil {
		return nil, err
	}
	if ok {
		return data.NewList(s, v), nil
	}
	return nil, r.errorf(ErrPrefixedNotPaired, s)
}

func (r *parser) list() (data.Value, error) {
	res := data.Vector{}
	var sawDotAt = -1
	for pos := 0; ; pos++ {
		t, err := r.nextToken()
		if err != nil {
			return nil, err
		}
		if t == nil {
			break
		}
		switch t.Type() {
		case lex.Dot:
			if pos == 0 || sawDotAt != -1 {
				return nil, r.error(ErrInvalidListSyntax)
			}
			sawDotAt = pos
		case lex.ListEnd:
			if sawDotAt == -1 {
				return data.NewList(res...), nil
			} else if sawDotAt != len(res)-1 {
				return nil, r.error(ErrInvalidListSyntax)
			}
			return makeDottedList(res...), nil
		default:
			v, err := r.value(t)
			if err != nil {
				return nil, err
			}
			res = append(res, v)
		}
	}
	return nil, r.error(ErrListNotClosed)
}

func (r *parser) vector() (data.Value, error) {
	return r.nonDotted(lex.VectorEnd)
}

func (r *parser) object() (data.Value, error) {
	v, err := r.nonDotted(lex.ObjectEnd)
	if err != nil {
		return nil, err
	}
	res, err := data.ValuesToObject(v...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *parser) nonDotted(endToken lex.TokenType) (data.Vector, error) {
	res := data.Vector{}
	for {
		t, err := r.nextToken()
		if err != nil {
			return nil, err
		}
		if t == nil {
			return nil, r.error(collectionErrors[endToken])
		}
		switch t.Type() {
		case endToken:
			return res, nil
		default:
			v, err := r.value(t)
			if err != nil {
				return nil, err
			}
			res = append(res, v)
		}
	}
}

func (r *parser) error(text string) error {
	return errors.New(text)
}

func (r *parser) errorf(text string, a ...any) error {
	return fmt.Errorf(text, a...)
}

func (r *parser) keyword() data.Value {
	n := r.token.Value().(data.String)
	return data.Keyword(n[1:])
}

func (r *parser) identifier() (data.Value, error) {
	n := r.token.Value().(data.String)
	if v, ok := specialNames[n]; ok {
		return v, nil
	}

	sym, err := data.ParseSymbol(n)
	if err != nil {
		return nil, err
	}
	return sym, nil
}

func makeDottedList(v ...data.Value) data.Value {
	l := len(v)
	if res, ok := v[l-1].(*data.List); ok {
		for i := l - 2; i >= 0; i-- {
			res = res.Prepend(v[i]).(*data.List)
		}
		return res
	}
	var res = data.NewCons(v[l-2], v[l-1])
	for i := l - 3; i >= 0; i-- {
		res = data.NewCons(v[i], res)
	}
	return res
}
