package read

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type (
	// TokenType is an opaque type for lexer tokens
	TokenType int

	// Token is a lexer value
	Token struct {
		Type  TokenType
		Value data.Value
	}

	tokenizer func([]string) *Token

	matchEntry struct {
		pattern  *regexp.Regexp
		function tokenizer
	}

	matchEntries []matchEntry
)

// Token Types
const (
	Error TokenType = iota
	Identifier
	Dot
	String
	Number
	ListStart
	ListEnd
	VectorStart
	VectorEnd
	MapStart
	MapEnd
	QuoteMarker
	SyntaxMarker
	UnquoteMarker
	SpliceMarker
	PatternMarker
	Whitespace
	Comment
	endOfFile
)

// Error messages
const (
	ErrStringNotTerminated = "string has no closing quote"
	ErrUnexpectedCharacter = "unexpected character: %s"

	errUnmatchedState = "unmatched lexing state"
)

const (
	structure  = `(){}\[\]\s\"`
	prefixChar = "`,~@"
	idStart    = "[^" + structure + prefixChar + "]"
	idCont     = "[^" + structure + "]"
	id         = idStart + idCont + "*"
	numTail    = idStart + "*"
)

var (
	escaped    = regexp.MustCompile(`\\\\|\\"|\\[^\\"]`)
	escapedMap = map[string]string{
		`\n`: "\n",
		`\s`: " ",
		`\t`: "\t",
		`\f`: "\f",
		`\b`: "\b",
		`\r`: "\r",
	}

	matchers = matchEntries{
		pattern(`$`, endState(endOfFile)),
		pattern(`;[^\n]*([\n]|$)`, tokenState(Comment)),
		pattern(`\s+`, tokenState(Whitespace)),
		pattern(`\(`, tokenState(ListStart)),
		pattern(`\[`, tokenState(VectorStart)),
		pattern(`{`, tokenState(MapStart)),
		pattern(`\)`, tokenState(ListEnd)),
		pattern(`]`, tokenState(VectorEnd)),
		pattern(`}`, tokenState(MapEnd)),
		pattern(`'`, tokenState(QuoteMarker)),
		pattern("`", tokenState(SyntaxMarker)),
		pattern(`,@`, tokenState(SpliceMarker)),
		pattern(`,`, tokenState(UnquoteMarker)),
		pattern(`~`, tokenState(PatternMarker)),

		pattern(`(")(?P<s>(\\\\|\\"|\\[^\\"]|[^"\\])*)("?)`, stringState),

		pattern(`[+-]?[1-9]\d*/[1-9]\d*`+numTail, ratioState),
		pattern(`[+-]?(0|[1-9]\d*)\.\d+([eE][+-]?\d+)?`+numTail, floatState),
		pattern(`[+-]?(0|[1-9]\d*)(\.\d+)?[eE][+-]?\d+`+numTail, floatState),
		pattern(`[+-]?0[bB]\d+`+numTail, integerState),
		pattern(`[+-]?0[xX][\dA-Fa-f]+`+numTail, integerState),
		pattern(`[+-]?0\d*`+numTail, integerState),
		pattern(`[+-]?[1-9]\d*`+numTail, integerState),

		pattern(id, identifierState),

		pattern(`.`, errorState),
	}
)

func pattern(p string, s tokenizer) matchEntry {
	return matchEntry{
		pattern:  regexp.MustCompile("^" + p),
		function: s,
	}
}

// Scan creates a new lexer Sequence
func Scan(src data.String) data.Sequence {
	var resolver sequence.LazyResolver
	s := string(src)

	resolver = func() (data.Value, data.Sequence, bool) {
		if t, rs := matchToken(s); t.Type != endOfFile {
			s = rs
			return t, sequence.NewLazy(resolver), true
		}
		return data.Nil, data.EmptyList, false
	}

	l := sequence.NewLazy(resolver)
	return sequence.Filter(l, notWhitespace)
}

var notWhitespace = data.Applicative(func(args ...data.Value) data.Value {
	t := args[0].(*Token)
	return data.Bool(t.Type != Whitespace && t.Type != Comment)
}, 1)

func matchToken(src string) (*Token, string) {
	for _, s := range matchers {
		if sm := s.pattern.FindStringSubmatch(src); sm != nil {
			return s.function(sm), src[len(sm[0]):]
		}
	}
	// Programmer error
	panic(errors.New(errUnmatchedState))
}

// Equal compares this Token to another for equality
func (t *Token) Equal(v data.Value) bool {
	if v, ok := v.(*Token); ok {
		return t.Type == v.Type && t.Value.Equal(v.Value)
	}
	return false
}

// String converts this Value into a string
func (t *Token) String() string {
	return data.NewVector(data.Integer(t.Type), t.Value).String()
}

func makeToken(t TokenType, v data.Value) *Token {
	return &Token{
		Type:  t,
		Value: v,
	}
}

func tokenState(t TokenType) tokenizer {
	return func(sm []string) *Token {
		return makeToken(t, data.String(sm[0]))
	}
}

func endState(t TokenType) tokenizer {
	return func(_ []string) *Token {
		return makeToken(t, nil)
	}
}

func unescape(s string) string {
	r := escaped.ReplaceAllStringFunc(s, func(e string) string {
		if r, ok := escapedMap[e]; ok {
			return r
		}
		return e[1:]
	})
	return r
}

func stringState(sm []string) *Token {
	if len(sm[4]) == 0 {
		return makeToken(Error, data.String(ErrStringNotTerminated))
	}
	s := unescape(sm[2])
	return makeToken(String, data.String(s))
}

func ratioState(sm []string) *Token {
	return tokenizeNumber(data.ParseRatio(sm[0]))
}

func floatState(sm []string) *Token {
	return tokenizeNumber(data.ParseFloat(sm[0]))
}

func integerState(sm []string) *Token {
	return tokenizeNumber(data.ParseInteger(sm[0]))
}

func tokenizeNumber(res data.Number, err error) *Token {
	if err != nil {
		return makeToken(Error, data.String(err.Error()))
	}
	return makeToken(Number, res)
}

func identifierState(sm []string) *Token {
	s := data.String(sm[0])
	if s == "." {
		return makeToken(Dot, s)
	}
	return makeToken(Identifier, s)
}

func errorState(sm []string) *Token {
	err := fmt.Errorf(ErrUnexpectedCharacter, sm[0])
	return makeToken(Error, data.String(err.Error()))
}
