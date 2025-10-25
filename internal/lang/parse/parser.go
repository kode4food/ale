package parse

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/lang/lex"
)

type (
	// parser is a stateful iteration interface for a Token stream that is
	// piloted by the FromString function and exposed as a LazySequence
	parser struct {
		ns       env.Namespace
		tokenize Tokenizer
		seq      data.Sequence
		token    *lex.Token
	}

	handler func(*parser, *lex.Token) (ale.Value, error)
)

var (
	ErrPrefixedNotPaired  = errors.New("end of file before completing")
	ErrUnexpectedDot      = errors.New("unexpected dot")
	ErrInvalidListSyntax  = errors.New("invalid list syntax")
	ErrListNotClosed      = errors.New("list not closed")
	ErrUnmatchedListEnd   = errors.New("unmatched list end")
	ErrVectorNotClosed    = errors.New("vector not closed")
	ErrUnmatchedVectorEnd = errors.New("unmatched vector end")
	ErrObjectNotClosed    = errors.New("object not closed")
	ErrUnmatchedObjectEnd = errors.New("unmatched object end")
)

var (
	quoteSym    = env.RootSymbol("quote")
	syntaxSym   = env.RootSymbol("syntax-quote")
	unquoteSym  = env.RootSymbol("unquote")
	splicingSym = env.RootSymbol("unquote-splicing")
	bytesSym    = env.RootSymbol("bytes")

	specialNames = map[data.String]ale.Value{
		lang.TrueLiteral:  data.True,
		lang.FalseLiteral: data.False,
	}

	collectionErrors = map[lex.TokenType]error{
		lex.VectorEnd: ErrVectorNotClosed,
		lex.ObjectEnd: ErrObjectNotClosed,
	}

	handlers     [lex.EOF + 1]handler
	handlersOnce sync.Once
)

func (p *parser) nextValue() (ale.Value, bool, error) {
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
		return nil, errors.New(data.ToString(p.token.Value()))
	}
	return p.token, nil
}

func (p *parser) maybeWrap(err error) error {
	if t := p.token; t != nil {
		return t.WrapError(err)
	}
	return err
}

func (p *parser) value(t *lex.Token) (ale.Value, error) {
	handlers := getValueHandlers()
	if handler := handlers[t.Type()]; handler != nil {
		return handler(p, t)
	}
	return t.Value(), nil
}

func (p *parser) prefixed(s data.Symbol) (ale.Value, error) {
	v, ok, err := p.nextValue()
	if err != nil {
		return nil, err
	}
	if ok {
		return data.NewList(s, v), nil
	}
	return nil, p.errorf(ErrPrefixedNotPaired, s)
}

func (p *parser) list() (ale.Value, error) {
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

func (p *parser) bytes() (ale.Value, error) {
	v, err := p.nonDotted(lex.VectorEnd)
	if err != nil {
		return nil, err
	}
	if res, err := data.ValuesToBytes(v...); err == nil {
		return res, nil
	}
	return data.NewList(append(data.Vector{bytesSym}, v...)...), nil
}

func (p *parser) vector() (ale.Value, error) {
	return p.nonDotted(lex.VectorEnd)
}

func (p *parser) object() (ale.Value, error) {
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

func (p *parser) error(err error) error {
	return err
}

func (p *parser) errorf(err error, a ...any) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprint(a...))
}

func (p *parser) keyword() (ale.Value, error) {
	n := p.token.Value().(data.String)
	return data.Keyword(n[1:]), nil
}

func (p *parser) identifier() (ale.Value, error) {
	n := p.token.Value().(data.String)
	if v, ok := specialNames[n]; ok {
		return v, nil
	}
	return data.ParseSymbol(n)
}

func makeDottedList(vals ...ale.Value) ale.Value {
	l := len(vals)
	if res, ok := vals[l-1].(*data.List); ok {
		for i := l - 2; i >= 0; i-- {
			res = res.Prepend(vals[i]).(*data.List)
		}
		return res
	}
	var res = data.NewCons(vals[l-2], vals[l-1])
	for i := l - 3; i >= 0; i-- {
		res = data.NewCons(vals[i], res)
	}
	return res
}

func getValueHandlers() *[lex.EOF + 1]handler {
	handlersOnce.Do(func() {
		handlers[lex.QuoteMarker] = makePrefixedHandler(quoteSym)
		handlers[lex.SyntaxMarker] = makePrefixedHandler(syntaxSym)
		handlers[lex.UnquoteMarker] = makePrefixedHandler(unquoteSym)
		handlers[lex.SpliceMarker] = makePrefixedHandler(splicingSym)
		handlers[lex.ListStart] = listStartHandler
		handlers[lex.BytesStart] = makeMethodHandler((*parser).bytes)
		handlers[lex.VectorStart] = makeMethodHandler((*parser).vector)
		handlers[lex.ObjectStart] = makeMethodHandler((*parser).object)
		handlers[lex.Keyword] = makeMethodHandler((*parser).keyword)
		handlers[lex.Identifier] = makeMethodHandler((*parser).identifier)
		handlers[lex.ListEnd] = makeErrorHandler(ErrUnmatchedListEnd)
		handlers[lex.VectorEnd] = makeErrorHandler(ErrUnmatchedVectorEnd)
		handlers[lex.ObjectEnd] = makeErrorHandler(ErrUnmatchedObjectEnd)
		handlers[lex.Dot] = makeErrorHandler(ErrUnexpectedDot)
	})
	return &handlers
}

func listStartHandler(p *parser, t *lex.Token) (ale.Value, error) {
	res, err := p.list()
	if err != nil {
		return nil, err
	}
	return p.processInclude(res)
}

func makePrefixedHandler(s data.Symbol) handler {
	return func(p *parser, t *lex.Token) (ale.Value, error) {
		return p.prefixed(s)
	}
}

func makeMethodHandler(method func(*parser) (ale.Value, error)) handler {
	return func(p *parser, t *lex.Token) (ale.Value, error) {
		return method(p)
	}
}

func makeErrorHandler(err error) handler {
	return func(p *parser, t *lex.Token) (ale.Value, error) {
		return nil, err
	}
}
