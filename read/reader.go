package read

import (
	"fmt"
	"regexp"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
)

// reader is a stateful iteration interface for a token stream
type reader struct {
	seq data.Sequence
}

// Error messages
const (
	PrefixedNotPaired  = "end of file reached before completing %s"
	ListNotClosed      = "end of file reached with open list"
	UnmatchedListEnd   = "encountered ')' with no open list"
	VectorNotClosed    = "end of file reached with open vector"
	UnmatchedVectorEnd = "encountered ']' with no open vector"
	MapNotClosed       = "end of file reached with open map"
	UnmatchedMapEnd    = "encountered '}' with no open map"
	MapNotPaired       = "map does not contain an even number of elements"
)

var (
	keywordIdentifier = regexp.MustCompile(`^:[^(){}\[\]\s,]+`)

	quoteSym    = namespace.RootSymbol("quote")
	syntaxSym   = namespace.RootSymbol("syntax-quote")
	unquoteSym  = namespace.RootSymbol("unquote")
	splicingSym = namespace.RootSymbol("unquote-splicing")

	specialNames = map[data.Name]data.Value{
		"true":  data.True,
		"false": data.False,
		"nil":   data.Nil,
	}
)

func newReader(lexer data.Sequence) *reader {
	return &reader{
		seq: lexer,
	}
}

func (r *reader) nextToken() *Token {
	s := r.seq
	if !s.IsSequence() {
		return nil
	}
	f := s.First()
	r.seq = s.Rest()
	return f.(*Token)
}

func (r *reader) nextValue() (data.Value, bool) {
	if t := r.nextToken(); t != nil {
		return r.value(t), true
	}
	return nil, false
}

func (r *reader) value(t *Token) data.Value {
	switch t.Type {
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
	case MapStart:
		return r.associative()
	case Identifier:
		return readIdentifier(t)
	case ListEnd:
		panic(fmt.Errorf(UnmatchedListEnd))
	case VectorEnd:
		panic(fmt.Errorf(UnmatchedVectorEnd))
	case MapEnd:
		panic(fmt.Errorf(UnmatchedMapEnd))
	default:
		return t.Value
	}
}

func (r *reader) prefixed(s data.Symbol) data.Value {
	if v, ok := r.nextValue(); ok {
		return data.NewList(s, v)
	}
	panic(fmt.Errorf(PrefixedNotPaired, s))
}

func (r *reader) list() data.Value {
	var handle func(t *Token) *data.List
	var rest func() *data.List

	handle = func(t *Token) *data.List {
		switch t.Type {
		case ListEnd:
			return data.EmptyList
		default:
			v := r.value(t)
			l := rest()
			return l.Prepend(v).(*data.List)
		}
	}

	rest = func() *data.List {
		if t := r.nextToken(); t != nil {
			return handle(t)
		}
		panic(fmt.Errorf(ListNotClosed))
	}

	return rest()
}

func (r *reader) vector() data.Value {
	res := make(data.Vector, 0)

	for {
		if t := r.nextToken(); t != nil {
			switch t.Type {
			case VectorEnd:
				return res
			default:
				e := r.value(t)
				res = append(res, e)
			}
		} else {
			panic(fmt.Errorf(VectorNotClosed))
		}
	}
}

func (r *reader) associative() data.Value {
	res := make([]data.Vector, 0)
	mp := make(data.Vector, 2)

	for idx := 0; ; idx++ {
		if t := r.nextToken(); t != nil {
			switch t.Type {
			case MapEnd:
				if idx%2 == 0 {
					return data.NewAssociative(res...)
				}
				panic(fmt.Errorf(MapNotPaired))
			default:
				e := r.value(t)
				if idx%2 == 0 {
					mp[0] = e
				} else {
					mp[1] = e
					res = append(res, mp)
					mp = make(data.Vector, 2)
				}
			}
		} else {
			panic(fmt.Errorf(MapNotClosed))
		}
	}
}

func readIdentifier(t *Token) data.Value {
	n := data.Name(t.Value.(data.String))
	if v, ok := specialNames[n]; ok {
		return v
	}

	s := string(n)
	if keywordIdentifier.MatchString(s) {
		return data.Keyword(n[1:])
	}
	return data.ParseSymbol(n)
}
