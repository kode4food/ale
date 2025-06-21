package parse

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

// parser is a stateful iteration interface for a Token stream that is piloted
// by the FromString function and exposed as a LazySequence
type parser struct {
	ns       env.Namespace
	tokenize Tokenizer
	seq      data.Sequence
	token    *lex.Token
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
	bytesSym    = env.RootSymbol("bytes")

	specialNames = map[data.String]data.Value{
		data.TrueLiteral:  data.True,
		data.FalseLiteral: data.False,
	}

	collectionErrors = map[lex.TokenType]string{
		lex.VectorEnd: ErrVectorNotClosed,
		lex.ObjectEnd: ErrObjectNotClosed,
	}
)

func (p *parser) nextValue() (data.Value, bool, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, false, p.maybeWrap(err)
	}
	if t != nil {
		v, err := p.value(t)
		if err != nil {
			return nil, false, p.maybeWrap(err)

		}
		return v, true, nil
	}
	return nil, false, nil
}

func (p *parser) nextToken() (*lex.Token, error) {
	token, seq, ok := p.seq.Split()
	if !ok {
		return nil, nil
	}
	p.token = token.(*lex.Token)
	p.seq = seq
	if p.token.Type() == lex.Error {
		return nil, p.error(data.ToString(p.token.Value()))
	}
	return p.token, nil
}

func (p *parser) maybeWrap(err error) error {
	if t := p.token; t != nil {
		return t.WrapError(err)
	}
	return err
}

func (p *parser) value(t *lex.Token) (data.Value, error) {
	switch t.Type() {
	case lex.QuoteMarker:
		return p.prefixed(quoteSym)
	case lex.SyntaxMarker:
		return p.prefixed(syntaxSym)
	case lex.UnquoteMarker:
		return p.prefixed(unquoteSym)
	case lex.SpliceMarker:
		return p.prefixed(splicingSym)
	case lex.ListStart:
		return p.processInclude(p.list())
	case lex.BytesStart:
		return p.bytes()
	case lex.VectorStart:
		return p.vector()
	case lex.ObjectStart:
		return p.object()
	case lex.Keyword:
		return p.keyword(), nil
	case lex.Identifier:
		return p.identifier()
	case lex.ListEnd:
		return nil, p.error(ErrUnmatchedListEnd)
	case lex.VectorEnd:
		return nil, p.error(ErrUnmatchedVectorEnd)
	case lex.ObjectEnd:
		return nil, p.error(ErrUnmatchedObjectEnd)
	case lex.Dot:
		return nil, p.error(ErrUnexpectedDot)
	default:
		return t.Value(), nil
	}
}

func (p *parser) prefixed(s data.Symbol) (data.Value, error) {
	v, ok, err := p.nextValue()
	if err != nil {
		return nil, err
	}
	if ok {
		return data.NewList(s, v), nil
	}
	return nil, p.errorf(ErrPrefixedNotPaired, s)
}

func (p *parser) list() (data.Value, error) {
	res := data.Vector{}
	var sawDotAt = -1
	for pos := 0; ; pos++ {
		t, err := p.nextToken()
		if err != nil {
			return nil, err
		}
		if t == nil {
			break
		}
		switch t.Type() {
		case lex.Dot:
			if pos == 0 || sawDotAt != -1 {
				return nil, p.error(ErrInvalidListSyntax)
			}
			sawDotAt = pos
		case lex.ListEnd:
			if sawDotAt == -1 {
				return data.NewList(res...), nil
			} else if sawDotAt != len(res)-1 {
				return nil, p.error(ErrInvalidListSyntax)
			}
			return makeDottedList(res...), nil
		default:
			v, err := p.value(t)
			if err != nil {
				return nil, err
			}
			res = append(res, v)
		}
	}
	return nil, p.error(ErrListNotClosed)
}

func (p *parser) bytes() (data.Value, error) {
	v, err := p.nonDotted(lex.VectorEnd)
	if err != nil {
		return nil, err
	}
	if res, err := data.ValuesToBytes(v...); err == nil {
		return res, nil
	}
	return data.NewList(append(data.Vector{bytesSym}, v...)...), nil
}

func (p *parser) vector() (data.Value, error) {
	return p.nonDotted(lex.VectorEnd)
}

func (p *parser) object() (data.Value, error) {
	v, err := p.nonDotted(lex.ObjectEnd)
	if err != nil {
		return nil, err
	}
	return data.ValuesToObject(v...)
}

func (p *parser) nonDotted(endToken lex.TokenType) (data.Vector, error) {
	res := data.Vector{}
	for {
		t, err := p.nextToken()
		if err != nil {
			return nil, err
		}
		if t == nil {
			return nil, p.error(collectionErrors[endToken])
		}
		switch t.Type() {
		case endToken:
			return res, nil
		default:
			v, err := p.value(t)
			if err != nil {
				return nil, err
			}
			res = append(res, v)
		}
	}
}

func (p *parser) error(text string) error {
	return errors.New(text)
}

func (p *parser) errorf(text string, a ...any) error {
	return fmt.Errorf(text, a...)
}

func (p *parser) keyword() data.Value {
	n := p.token.Value().(data.String)
	return data.Keyword(n[1:])
}

func (p *parser) identifier() (data.Value, error) {
	n := p.token.Value().(data.String)
	if v, ok := specialNames[n]; ok {
		return v, nil
	}
	return data.ParseSymbol(n)
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
