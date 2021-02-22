package read

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type (
	tokenizer func([]string) *Token

	matchEntry struct {
		pattern  *regexp.Regexp
		function tokenizer
	}

	matchEntries []matchEntry
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

// Scan creates a new lexer Sequence
func Scan(src data.String) data.Sequence {
	var resolver sequence.LazyResolver
	var line, column int
	input := string(src)

	resolver = func() (data.Value, data.Sequence, bool) {
		if t, rest := matchToken(input); t.Type() != endOfFile {
			t := t.WithLocation(line, column)

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

func pattern(p string, s tokenizer) matchEntry {
	return matchEntry{
		pattern:  regexp.MustCompile("^" + p),
		function: s,
	}
}

func matchToken(src string) (*Token, string) {
	for _, s := range matchers {
		if sm := s.pattern.FindStringSubmatch(src); sm != nil {
			return s.function(sm), src[len(sm[0]):]
		}
	}
	// Programmer error
	panic(errors.New(errUnmatchedState))
}

func tokenState(t TokenType) tokenizer {
	return func(sm []string) *Token {
		return MakeToken(t, data.String(sm[0]))
	}
}

func endState(t TokenType) tokenizer {
	return func(_ []string) *Token {
		return MakeToken(t, nil)
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
		return MakeToken(Error, data.String(ErrStringNotTerminated))
	}
	s := unescape(sm[2])
	return MakeToken(String, data.String(s))
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
		return MakeToken(Error, data.String(err.Error()))
	}
	return MakeToken(Number, res)
}

func identifierState(sm []string) *Token {
	s := data.String(sm[0])
	if s == "." {
		return MakeToken(Dot, s)
	}
	return MakeToken(Identifier, s)
}

func errorState(sm []string) *Token {
	err := fmt.Errorf(ErrUnexpectedCharacter, sm[0])
	return MakeToken(Error, data.String(err.Error()))
}
