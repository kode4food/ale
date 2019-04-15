package read

import (
	"fmt"
	"regexp"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/stdlib"
)

const (
	// PrefixedNotPaired is thrown when a Quote is not completed
	PrefixedNotPaired = "end of file reached before completing %s"

	// ListNotClosed is thrown when EOF is reached inside a List
	ListNotClosed = "end of file reached with open list"

	// UnmatchedListEnd is thrown if a list is ended without being started
	UnmatchedListEnd = "encountered ')' with no open list"

	// VectorNotClosed is thrown when EOF is reached inside a Vector
	VectorNotClosed = "end of file reached with open vector"

	// UnmatchedVectorEnd is thrown if a vector is ended without being started
	UnmatchedVectorEnd = "encountered ']' with no open vector"

	// MapNotClosed is thrown when EOF is reached inside a Manager
	MapNotClosed = "end of file reached with open map"

	// UnmatchedMapEnd is thrown if a list is ended without being started
	UnmatchedMapEnd = "encountered '}' with no open map"

	// MapNotPaired is thrown if a map doesn't have an even number of elements
	MapNotPaired = "map does not contain an even number of elements"
)

type reader struct {
	iter *stdlib.Iterator
}

var (
	keywordIdentifier = regexp.MustCompile(`^:[^(){}\[\]\s,]+`)

	quoteSym    = namespace.RootSymbol("quote")
	syntaxSym   = namespace.RootSymbol("syntax-quote")
	unquoteSym  = namespace.RootSymbol("unquote")
	splicingSym = namespace.RootSymbol("unquote-splicing")

	specialNames = map[api.Name]api.Value{
		"true":  api.True,
		"false": api.False,
		"nil":   api.Nil,
	}
)

// FromString converts the raw source into unexpanded data structures
func FromString(src api.String) api.Sequence {
	l := Scan(src)
	return FromScanner(l)
}

// FromScanner returns a Lazy Sequence of scanned data structures
func FromScanner(lexer api.Sequence) api.Sequence {
	var res stdlib.LazyResolver
	r := newReader(lexer)

	res = func() (api.Value, api.Sequence, bool) {
		if f, ok := r.nextValue(); ok {
			return f, stdlib.NewLazySequence(res), true
		}
		return api.Nil, api.EmptyList, false
	}

	return stdlib.NewLazySequence(res)
}

func newReader(lexer api.Sequence) *reader {
	return &reader{
		iter: stdlib.Iterate(lexer),
	}
}

func (r *reader) nextToken() *Token {
	if t, ok := r.iter.Next(); ok {
		return t.(*Token)
	}
	return nil
}

func (r *reader) nextValue() (api.Value, bool) {
	if t := r.nextToken(); t != nil {
		return r.value(t), true
	}
	return nil, false
}

func (r *reader) value(t *Token) api.Value {
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

func (r *reader) prefixed(s api.Symbol) api.Value {
	if v, ok := r.nextValue(); ok {
		return api.NewList(s, v)
	}
	panic(fmt.Errorf(PrefixedNotPaired, s))
}

func (r *reader) list() api.Value {
	var handle func(t *Token) *api.List
	var rest func() *api.List

	handle = func(t *Token) *api.List {
		switch t.Type {
		case ListEnd:
			return api.EmptyList
		default:
			v := r.value(t)
			l := rest()
			return l.Prepend(v).(*api.List)
		}
	}

	rest = func() *api.List {
		if t := r.nextToken(); t != nil {
			return handle(t)
		}
		panic(fmt.Errorf(ListNotClosed))
	}

	return rest()
}

func (r *reader) vector() api.Value {
	res := make(api.Vector, 0)

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

func (r *reader) associative() api.Value {
	res := make([]api.Vector, 0)
	mp := make(api.Vector, 2)

	for idx := 0; ; idx++ {
		if t := r.nextToken(); t != nil {
			switch t.Type {
			case MapEnd:
				if idx%2 == 0 {
					return api.NewAssociative(res...)
				}
				panic(fmt.Errorf(MapNotPaired))
			default:
				e := r.value(t)
				if idx%2 == 0 {
					mp[0] = e
				} else {
					mp[1] = e
					res = append(res, mp)
					mp = make(api.Vector, 2)
				}
			}
		} else {
			panic(fmt.Errorf(MapNotClosed))
		}
	}
}

func readIdentifier(t *Token) api.Value {
	n := api.Name(t.Value.(api.String))
	if v, ok := specialNames[n]; ok {
		return v
	}

	s := string(n)
	if keywordIdentifier.MatchString(s) {
		return api.Keyword(n[1:])
	}
	return api.ParseSymbol(n)
}
