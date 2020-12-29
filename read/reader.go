package read

import (
	"fmt"
	"regexp"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
)

// reader is a stateful iteration interface for a token stream
type reader struct {
	seq data.Sequence
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
	keywordIdentifier = regexp.MustCompile(`^:[^(){}\[\]\s,]+`)

	quoteSym    = env.RootSymbol("quote")
	syntaxSym   = env.RootSymbol("syntax-quote")
	unquoteSym  = env.RootSymbol("unquote")
	splicingSym = env.RootSymbol("unquote-splicing")
	patternSym  = env.RootSymbol("pattern")

	specialNames = map[data.String]data.Value{
		data.TrueLiteral:  data.True,
		data.FalseLiteral: data.False,
	}

	collectionErrors = map[TokenType]string{
		VectorEnd: ErrVectorNotClosed,
		MapEnd:    ErrMapNotClosed,
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
	if f := s.First().(*Token); f.Type == Error {
		panic(fmt.Errorf("%s at line %d:%d", f.Value.String(), f.Line, f.Column))
	} else {
		r.seq = s.Rest()
		return f
	}
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
		return r.list(t)
	case VectorStart:
		return r.vector(t)
	case MapStart:
		return r.object(t)
	case Identifier:
		return readIdentifier(t)
	case ListEnd:
		panic(fmt.Errorf("%s at line %d:%d", ErrUnmatchedListEnd, t.Line, t.Column))
	case VectorEnd:
		panic(fmt.Errorf("%s at line %d:%d", ErrUnmatchedVectorEnd, t.Line, t.Column))
	case MapEnd:
		panic(fmt.Errorf("%s at line %d:%d", ErrUnmatchedMapEnd, t.Line, t.Column))
	case Dot:
		panic(fmt.Errorf("%s at line %d:%d", ErrUnexpectedDot, t.Line, t.Column))
	default:
		return t.Value
	}
}

func (r *reader) prefixed(s data.Symbol) data.Value {
	if v, ok := r.nextValue(); ok {
		return data.NewList(s, v)
	}
	panic(fmt.Errorf(ErrPrefixedNotPaired, s))
}

func (r *reader) list(firstToken *Token) data.Value {
	res := data.Values{}
	var sawDotAt = -1
	for i := 0; ; i++ {
		if t := r.nextToken(); t != nil {
			switch t.Type {
			case Dot:
				if i == 0 || sawDotAt != -1 {
					panic(fmt.Errorf("%s at line %d:%d", ErrInvalidListSyntax, t.Line, t.Column))
				}
				sawDotAt = i
			case ListEnd:
				if sawDotAt == -1 {
					return data.NewList(res...)
				} else if sawDotAt != len(res)-1 {
					panic(fmt.Errorf("%s at line %d:%d", ErrInvalidListSyntax, t.Line, t.Column))
				}
				return makeDottedList(res...)
			default:
				res = append(res, r.value(t))
			}
		} else {
			panic(fmt.Errorf("%s starting at line %d:%d", ErrListNotClosed, firstToken.Line, firstToken.Column))
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

func (r *reader) vector(firstToken *Token) data.Value {
	v := r.readNonDotted(firstToken, VectorEnd)
	return data.NewVector(v...)
}

func (r *reader) object(firstToken *Token) data.Value {
	v := r.readNonDotted(firstToken, MapEnd)
	return data.ValuesToObject(v...)
}

func (r *reader) readNonDotted(firstToken *Token, endToken TokenType) data.Values {
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
			panic(fmt.Errorf("%s at line %d:%d", err, firstToken.Line, firstToken.Column))
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
