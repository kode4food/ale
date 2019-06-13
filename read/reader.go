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
)

var (
	keywordIdentifier = regexp.MustCompile(`^:[^(){}\[\]\s,]+`)

	quoteSym    = namespace.RootSymbol("quote")
	syntaxSym   = namespace.RootSymbol("syntax-quote")
	unquoteSym  = namespace.RootSymbol("unquote")
	splicingSym = namespace.RootSymbol("unquote-splicing")
	patternSym  = namespace.RootSymbol("pattern")

	specialNames = map[data.String]data.Value{
		"true":  data.True,
		"false": data.False,
		"null":  data.Null,
	}

	collectionErrors = map[TokenType]string{
		ListEnd:   ListNotClosed,
		VectorEnd: VectorNotClosed,
		MapEnd:    MapNotClosed,
	}
)

func newReader(lexer data.Sequence) *reader {
	return &reader{
		seq: lexer,
	}
}

func (r *reader) nextToken() *Token {
	s := r.seq
	if s.IsEmpty() {
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
	case PatternMarker:
		return r.prefixed(patternSym)
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
	elems := r.readCollection(ListEnd)
	return data.NewList(elems...)
}

func (r *reader) vector() data.Value {
	elems := r.readCollection(VectorEnd)
	return data.NewVector(elems...)
}

func (r *reader) associative() data.Value {
	elems := r.readCollection(MapEnd)
	return data.NewAssociative(elems...)
}

func (r *reader) readCollection(endToken TokenType) data.Values {
	res := data.Values{}
	for {
		if t := r.nextToken(); t != nil {
			switch t.Type {
			case endToken:
				return res
			default:
				e := r.value(t)
				res = append(res, e)
			}
		} else {
			err := collectionErrors[endToken]
			panic(fmt.Errorf(err))
		}
	}
}

func readIdentifier(t *Token) data.Value {
	n := t.Value.(data.String)
	if v, ok := specialNames[n]; ok {
		return v
	}

	s := string(n)
	if keywordIdentifier.MatchString(s) {
		return data.Keyword(n[1:])
	}
	return data.ParseSymbol(n)
}
