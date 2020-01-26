package read

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/namespace"
)

// reader is a stateful iteration interface for a token stream
type reader struct {
	seq data.Sequence
}

// Error messages
const (
	PrefixedNotPaired  = "end of file reached before completing %s"
	UnexpectedDot      = "encountered '.' with no open list"
	InvalidListSyntax  = "invalid list syntax"
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
		data.TrueLiteral:  data.True,
		data.FalseLiteral: data.False,
	}

	collectionErrors = map[TokenType]string{
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
		return r.object()
	case Identifier:
		return readIdentifier(t)
	case ListEnd:
		panic(errors.New(UnmatchedListEnd))
	case VectorEnd:
		panic(errors.New(UnmatchedVectorEnd))
	case MapEnd:
		panic(errors.New(UnmatchedMapEnd))
	case Dot:
		panic(errors.New(UnexpectedDot))
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
	res := data.Values{}
	var sawDotAt = -1
	for i := 0; ; i++ {
		if t := r.nextToken(); t != nil {
			switch t.Type {
			case Dot:
				if i == 0 || sawDotAt != -1 {
					panic(errors.New(InvalidListSyntax))
				}
				sawDotAt = i
			case ListEnd:
				if sawDotAt == -1 {
					return data.NewList(res...)
				} else if sawDotAt != len(res)-1 {
					panic(errors.New(InvalidListSyntax))
				}
				return makeDottedList(res...)
			default:
				res = append(res, r.value(t))
			}
		} else {
			panic(errors.New(ListNotClosed))
		}
	}
}

func makeDottedList(v ...data.Value) data.Value {
	l := len(v)
	var res = data.NewCons(v[l-2], v[l-1])
	for i := l - 3; i >= 0; i-- {
		res = data.NewCons(v[i], res)
	}
	return res
}

func (r *reader) vector() data.Value {
	v := r.readNonDotted(VectorEnd)
	return data.NewVector(v...)
}

func (r *reader) object() data.Value {
	v := r.readNonDotted(MapEnd)
	return data.NewObject(v...)
}

func (r *reader) readNonDotted(endToken TokenType) data.Values {
	res := data.Values{}
	for {
		if t := r.nextToken(); t != nil {
			switch t.Type {
			case endToken:
				return res
			default:
				res = append(res, r.value(t))
			}
		} else {
			err := collectionErrors[endToken]
			panic(errors.New(err))
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
