package read

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
)

// reader is a stateful iteration interface for a Token stream that is piloted
// by the FromScanner function and exposed as a LazySequence
type reader struct {
	seq   data.Sequence
	token *Token
}

// Error messages
const (
	ErrPrefixedNotPaired  = "end of file reached before completing %s"
	ErrUnexpectedDot      = "encountered '.' with no open list"
	ErrInvalidListSyntax  = "invalid list syntax"
	ErrListNotClosed      = "end of file reached with open list"
	ErrUnmatchedListEnd   = "encountered ')' with no open list"
	ErrVectorNotClosed    = "end of file reached with open vector"
	ErrUnmatchedVectorEnd = "encountered ']' with no open vector"
	ErrMapNotClosed       = "end of file reached with open map"
	ErrUnmatchedMapEnd    = "encountered '}' with no open map"
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

	collectionErrors = map[TokenType]string{
		VectorEnd: ErrVectorNotClosed,
		ObjectEnd: ErrMapNotClosed,
	}
)

func newReader(lexer data.Sequence) *reader {
	return &reader{
		seq: lexer,
	}
}

func (r *reader) nextValue() (data.Value, bool) {
	if t := r.nextToken(); t != nil {
		return r.value(t), true
	}
	return nil, false
}

func (r *reader) nextToken() *Token {
	token, seq, ok := r.seq.Split()
	if !ok {
		return nil
	}
	r.token = token.(*Token)
	r.seq = seq
	if r.token.Type() == Error {
		panic(r.error(r.token.Value().String()))
	}
	return r.token
}

func (r *reader) value(t *Token) data.Value {
	switch t.Type() {
	case QuoteMarker:
		return r.prefixed(quoteSym)
	case SyntaxMarker:
		return r.prefixed(syntaxSym)
	case UnquoteMarker:
		return r.prefixed(unquoteSym)
	case SpliceMarker:
		return r.prefixed(splicingSym)
	case ListStart:
		return r.list()
	case VectorStart:
		return r.vector()
	case ObjectStart:
		return r.object()
	case Keyword:
		return r.keyword()
	case Identifier:
		return r.identifier()
	case ListEnd:
		panic(r.error(ErrUnmatchedListEnd))
	case VectorEnd:
		panic(r.error(ErrUnmatchedVectorEnd))
	case ObjectEnd:
		panic(r.error(ErrUnmatchedMapEnd))
	case Dot:
		panic(r.error(ErrUnexpectedDot))
	default:
		return t.Value()
	}
}

func (r *reader) prefixed(s data.Symbol) data.Value {
	if v, ok := r.nextValue(); ok {
		return data.NewList(s, v)
	}
	panic(r.errorf(ErrPrefixedNotPaired, s))
}

func (r *reader) list() data.Value {
	res := data.Values{}
	var sawDotAt = -1
	for t, i := r.nextToken(), 0; t != nil; t, i = r.nextToken(), i+1 {
		switch t.Type() {
		case Dot:
			if i == 0 || sawDotAt != -1 {
				panic(r.error(ErrInvalidListSyntax))
			}
			sawDotAt = i
		case ListEnd:
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

func (r *reader) vector() data.Value {
	v := r.readNonDotted(VectorEnd)
	return data.NewVector(v...)
}

func (r *reader) object() data.Value {
	v := r.readNonDotted(ObjectEnd)
	res, err := data.ValuesToObject(v...)
	if err != nil {
		panic(r.maybeWrap(err))
	}
	return res
}

func (r *reader) readNonDotted(endToken TokenType) data.Values {
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

func (r *reader) maybeWrap(err error) error {
	if t := r.token; t != nil {
		return t.wrapError(err)
	}
	return err
}

func (r *reader) error(text string) error {
	return r.maybeWrap(errors.New(text))
}

func (r *reader) errorf(text string, a ...any) error {
	return r.maybeWrap(fmt.Errorf(text, a...))
}

func (r *reader) keyword() data.Value {
	n := r.token.Value().(data.String)
	return data.Keyword(n[1:])
}

func (r *reader) identifier() data.Value {
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
