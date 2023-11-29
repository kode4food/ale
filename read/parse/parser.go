package parse

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/read/lex"
)

// parser is a stateful iteration interface for a Token stream that is piloted
// by the FromLexer function and exposed as a LazySequence
type parser struct {
	seq   data.Sequence
	token *lex.Token
}

const (
	// ErrPrefixedNotPaired is raised when the parser encounters the end of the
	// stream without being able to completed a paired element, such as a quote
	ErrPrefixedNotPaired = "end of file reached before completing %s"

	// ErrUnexpectedDot is raised when the parser encounters a dot in the
	// stream when it isn't part of an open list
	ErrUnexpectedDot = "encountered '.' with no open list"

	// ErrInvalidListSyntax is raised when the parse encounters a misplaced dot
	// when parsing an open list
	ErrInvalidListSyntax = "invalid list syntax"

	// ErrListNotClosed is raised when the parser encounters the end of the
	// stream while an open list is still being parsed
	ErrListNotClosed = "end of file reached with open list"

	// ErrUnmatchedListEnd is raised when a list-end character is encountered
	// in the stream when no open list is being parsed
	ErrUnmatchedListEnd = "encountered ')' with no open list"

	// ErrVectorNotClosed is raised when the parser encounters the end of the
	// stream while an open vector is still being parsed
	ErrVectorNotClosed = "end of file reached with open vector"

	// ErrUnmatchedVectorEnd is raised when a vector-end character is
	// encountered in the stream when no open vector is being parsed
	ErrUnmatchedVectorEnd = "encountered ']' with no open vector"

	// ErrObjectNotClosed is raised when the parser encounters the end of the
	// stream while an open object is still being parsed
	ErrObjectNotClosed = "end of file reached with open object"

	// ErrUnmatchedObjectEnd is raised when an object-end character is
	// encountered in the stream when no open object is being parsed
	ErrUnmatchedObjectEnd = "encountered '}' with no open object"
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

func (r *parser) nextValue() (data.Value, bool) {
	if t := r.nextToken(); t != nil {
		return r.value(t), true
	}
	return nil, false
}

func (r *parser) nextToken() *lex.Token {
	token, seq, ok := r.seq.Split()
	if !ok {
		return nil
	}
	r.token = token.(*lex.Token)
	r.seq = seq
	if r.token.Type() == lex.Error {
		panic(r.error(data.ToString(r.token.Value())))
	}
	return r.token
}

func (r *parser) value(t *lex.Token) data.Value {
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
		return r.keyword()
	case lex.Identifier:
		return r.identifier()
	case lex.ListEnd:
		panic(r.error(ErrUnmatchedListEnd))
	case lex.VectorEnd:
		panic(r.error(ErrUnmatchedVectorEnd))
	case lex.ObjectEnd:
		panic(r.error(ErrUnmatchedObjectEnd))
	case lex.Dot:
		panic(r.error(ErrUnexpectedDot))
	default:
		return t.Value()
	}
}

func (r *parser) prefixed(s data.Symbol) data.Value {
	if v, ok := r.nextValue(); ok {
		return data.NewList(s, v)
	}
	panic(r.errorf(ErrPrefixedNotPaired, s))
}

func (r *parser) list() data.Value {
	res := data.Values{}
	var sawDotAt = -1
	for t, i := r.nextToken(), 0; t != nil; t, i = r.nextToken(), i+1 {
		switch t.Type() {
		case lex.Dot:
			if i == 0 || sawDotAt != -1 {
				panic(r.error(ErrInvalidListSyntax))
			}
			sawDotAt = i
		case lex.ListEnd:
			if sawDotAt == -1 {
				return data.NewList(res...)
			} else if sawDotAt != len(res)-1 {
				panic(r.error(ErrInvalidListSyntax))
			}
			return makeDottedList(res...)
		default:
			res = append(res, r.value(t))
		}
	}
	panic(r.error(ErrListNotClosed))
}

func (r *parser) vector() data.Value {
	v := r.nonDotted(lex.VectorEnd)
	return data.Vector(v)
}

func (r *parser) object() data.Value {
	v := r.nonDotted(lex.ObjectEnd)
	res, err := data.ValuesToObject(v...)
	if err != nil {
		panic(r.maybeWrap(err))
	}
	return res
}

func (r *parser) nonDotted(endToken lex.TokenType) data.Values {
	res := data.Values{}
	for {
		t := r.nextToken()
		if t == nil {
			panic(r.error(collectionErrors[endToken]))
		}
		switch t.Type() {
		case endToken:
			return res
		default:
			res = append(res, r.value(t))
		}
	}
}

func (r *parser) maybeWrap(err error) error {
	if t := r.token; t != nil {
		return t.WrapError(err)
	}
	return err
}

func (r *parser) error(text string) error {
	return r.maybeWrap(errors.New(text))
}

func (r *parser) errorf(text string, a ...any) error {
	return r.maybeWrap(fmt.Errorf(text, a...))
}

func (r *parser) keyword() data.Value {
	n := r.token.Value().(data.String)
	return data.Keyword(n[1:])
}

func (r *parser) identifier() data.Value {
	n := r.token.Value().(data.String)
	if v, ok := specialNames[n]; ok {
		return v
	}

	sym, err := data.ParseSymbol(n)
	if err != nil {
		panic(r.maybeWrap(err))
	}
	return sym
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
