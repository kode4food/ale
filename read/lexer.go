package read

import (
	"fmt"
	"regexp"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/stdlib"
)

// Error messages
const (
	UnmatchedState      = "unmatched lexing state"
	StringNotTerminated = "string has no closing quote"
)

// Token Types
const (
	Error TokenType = iota
	Identifier
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

const (
	idChar  = `[^(){}\[\]\s,'~@";]`
	id      = idChar + "*"
	numTail = id
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

		pattern(id, tokenState(Identifier)),

		pattern(`.`, endState(Error)),
	}
)

// Scan creates a new lexer Sequence
func Scan(src data.String) data.Sequence {
	var resolver stdlib.LazyResolver
	s := string(src)

	resolver = func() (data.Value, data.Sequence, bool) {
		if t, rs := matchToken(s); t.Type != endOfFile {
			s = rs
			return t, stdlib.NewLazySequence(resolver), true
		}
		return data.Null, data.EmptyList, false
	}

	l := stdlib.NewLazySequence(resolver)
	return stdlib.Filter(l, notWhitespace)
}

func notWhitespace(args ...data.Value) data.Value {
	t := args[0].(*Token)
	return data.Bool(t.Type != Whitespace && t.Type != Comment)
}

func matchToken(src string) (*Token, string) {
	for _, s := range matchers {
		if sm := s.pattern.FindStringSubmatch(src); sm != nil {
			return s.function(sm), src[len(sm[0]):]
		}
	}
	// Shouldn't happen because of the patterns that are defined,
	// but is here as a safety net
	panic(fmt.Errorf(UnmatchedState))
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
		panic(fmt.Errorf(StringNotTerminated))
	}
	s := unescape(sm[2])
	return makeToken(String, data.String(s))
}

func ratioState(sm []string) *Token {
	res := data.ParseRatio(sm[0])
	return makeToken(Number, res)
}

func floatState(sm []string) *Token {
	v := data.ParseFloat(sm[0])
	return makeToken(Number, v)
}

func integerState(sm []string) *Token {
	v := data.ParseInteger(sm[0])
	return makeToken(Number, v)
}

func pattern(p string, s tokenizer) matchEntry {
	return matchEntry{
		pattern:  regexp.MustCompile("^" + p),
		function: s,
	}
}
