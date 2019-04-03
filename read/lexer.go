package read

import (
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/stdlib"
)

const (
	// UnmatchedState is the error returned when the lexer state is invalid
	UnmatchedState = "unmatched lexing state"

	// StringNotTerminated is returned when the lexer sees a bad string
	StringNotTerminated = "string has no closing quote"
)

// Token Types
const (
	Error TokenType = iota
	Identifier
	String
	Number
	Ratio
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
		Value api.Value
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
	id      = idChar + "+"
	numTail = idChar + "*"
)

var (
	escaped    = regexp.MustCompile(`\\\\|\\"|\\[^\\"]`)
	escapedMap = map[string]string{
		`\\`: `\`,
		`\"`: `"`,
	}

	matchers = matchEntries{
		pattern(`$`, endState(endOfFile)),
		pattern(`;[^\n]*([\n]|$)`, tokenState(Comment)),
		pattern(`[\s,]+`, tokenState(Whitespace)),
		pattern(`\(`, tokenState(ListStart)),
		pattern(`\[`, tokenState(VectorStart)),
		pattern(`{`, tokenState(MapStart)),
		pattern(`\)`, tokenState(ListEnd)),
		pattern(`]`, tokenState(VectorEnd)),
		pattern(`}`, tokenState(MapEnd)),
		pattern(`'`, tokenState(QuoteMarker)),
		pattern("`", tokenState(SyntaxMarker)),
		pattern(`~@`, tokenState(SpliceMarker)),
		pattern(`~`, tokenState(UnquoteMarker)),

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
func Scan(src api.String) api.Sequence {
	var resolver stdlib.LazyResolver
	s := string(src)

	resolver = func() (api.Value, api.Sequence, bool) {
		if t, rs := matchToken(s); t.Type != endOfFile {
			s = rs
			return t, stdlib.NewLazySequence(resolver), true
		}
		return api.Nil, api.EmptyList, false
	}

	l := stdlib.NewLazySequence(resolver)
	return stdlib.Filter(l, notWhitespace)
}

func notWhitespace(args ...api.Value) api.Value {
	t := args[0].(*Token)
	return api.Bool(t.Type != Whitespace && t.Type != Comment)
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
	return api.Vector{api.Integer(t.Type), t.Value}.String()
}

func makeToken(t TokenType, v api.Value) *Token {
	return &Token{
		Type:  t,
		Value: v,
	}
}

func tokenState(t TokenType) tokenizer {
	return func(sm []string) *Token {
		return makeToken(t, api.String(sm[0]))
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
		return e
	})
	return r
}

func stringState(sm []string) *Token {
	if len(sm[4]) == 0 {
		panic(fmt.Errorf(StringNotTerminated))
	}
	s := unescape(sm[2])
	return makeToken(String, api.String(s))
}

func ratioState(sm []string) *Token {
	comp := strings.Split(sm[0], "/")
	num := int64(api.ParseInteger(api.String(comp[0])))
	den := int64(api.ParseInteger(api.String(comp[1])))
	res := api.Float(float64(num / den))
	return makeToken(Number, res)
}

func floatState(sm []string) *Token {
	v := api.ParseFloat(api.String(sm[0]))
	return makeToken(Number, v)
}

func integerState(sm []string) *Token {
	v := api.ParseInteger(api.String(sm[0]))
	return makeToken(Number, v)
}

func pattern(p string, s tokenizer) matchEntry {
	return matchEntry{
		pattern:  regexp.MustCompile("^" + p),
		function: s,
	}
}
