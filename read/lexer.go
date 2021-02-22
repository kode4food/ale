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
		Type   TokenType
		Value  data.Value
		Line   int
		Column int
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
	ObjectStart
	ObjectEnd
	QuoteMarker
	SyntaxMarker
	UnquoteMarker
	SpliceMarker
	PatternMarker
	Whitespace
	NewLine
	Comment
	endOfFile
)

// Error messages
const (
	ErrStringNotTerminated = "string has no closing quote"
	ErrUnexpectedCharacter = "unexpected character: %s"

	errTokenWrapped   = "%w (line %d, column %d)"
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
		pattern(`(\r\n|[\n\r])`, tokenState(NewLine)),
		pattern(`[\t\f ]+`, tokenState(Whitespace)),
		pattern(`\(`, tokenState(ListStart)),
		pattern(`\[`, tokenState(VectorStart)),
		pattern(`{`, tokenState(ObjectStart)),
		pattern(`\)`, tokenState(ListEnd)),
		pattern(`]`, tokenState(VectorEnd)),
		pattern(`}`, tokenState(ObjectEnd)),
		pattern(`'`, tokenState(QuoteMarker)),
		pattern("`", tokenState(SyntaxMarker)),
		pattern(`,@`, tokenState(SpliceMarker)),
		pattern(`,`, tokenState(UnquoteMarker)),
		pattern(`~`, tokenState(PatternMarker)),

		pattern(`(")(?P<s>(\\\\|\\"|\\[^\\"]|[^"\\])*)("?)`, stringState),

		pattern(`[+-]?(0|[1-9]\d*)/[1-9]\d*`+numTail, ratioState),
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
	var line, column int
	input := string(src)

	resolver = func() (data.Value, data.Sequence, bool) {
		if t, rest := matchToken(input); t.Type != endOfFile {
			t.Line = line
			t.Column = column

			if t.isNewLine() {
				line++
				column = 0
			} else {
				column += len(input) - len(rest)
			}

			input = rest
			return t, sequence.NewLazy(resolver), true
		}
		return data.Nil, data.EmptyList, false
	}

	res := sequence.NewLazy(resolver)
	return sequence.Filter(res, notWhitespace)
}

var notWhitespace = data.Applicative(func(args ...data.Value) data.Value {
	t := args[0].(*Token)
	return data.Bool(!t.isWhitespace())
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

func (t *Token) isNewLine() bool {
	return t.Type == Comment || t.Type == NewLine
}

func (t *Token) isWhitespace() bool {
	return t.Type == Comment || t.Type == NewLine || t.Type == Whitespace
}

func (t *Token) wrapError(e error) error {
	return fmt.Errorf(errTokenWrapped, e, t.Line+1, t.Column+1)
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
