package read

import (
	"fmt"
	"regexp"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/lang"
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

	dotRegex = regexp.MustCompile(`^` + lang.Dot + `$`)

	matchers = matchEntries{
		pattern(`$`, endState(endOfFile)),
		pattern(lang.Comment, tokenState(Comment)),
		pattern(lang.NewLine, tokenState(NewLine)),
		pattern(lang.Whitespace, tokenState(Whitespace)),
		pattern(lang.ListStart, tokenState(ListStart)),
		pattern(lang.VectorStart, tokenState(VectorStart)),
		pattern(lang.ObjectStart, tokenState(ObjectStart)),
		pattern(lang.ListEnd, tokenState(ListEnd)),
		pattern(lang.VectorEnd, tokenState(VectorEnd)),
		pattern(lang.ObjectEnd, tokenState(ObjectEnd)),
		pattern(lang.Quote, tokenState(QuoteMarker)),
		pattern(lang.SyntaxQuote, tokenState(SyntaxMarker)),
		pattern(lang.Splice, tokenState(SpliceMarker)),
		pattern(lang.Unquote, tokenState(UnquoteMarker)),
		pattern(lang.Pattern, tokenState(PatternMarker)),

		pattern(lang.String, stringState),

		pattern(lang.Ratio, ratioState),
		pattern(lang.Float, floatState),
		pattern(lang.Integer, integerState),

		pattern(lang.Keyword, tokenState(Keyword)),
		pattern(lang.Identifier, identifierState),

		pattern(lang.AnyChar, errorState),
	}
)

// Tokens creates a new Lexer Sequence of raw Tokens
func Tokens(src data.String) data.Sequence {
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
	return sequence.NewLazy(resolver)
}

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
	panic("unmatched lexing state")
}

func tokenState(t TokenType) tokenizer {
	return func(sm []string) *Token {
		return MakeToken(t, data.String(sm[0]))
	}
}

func endState(t TokenType) tokenizer {
	return func([]string) *Token {
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
	if dotRegex.MatchString(string(s)) {
		return MakeToken(Dot, s)
	}
	return MakeToken(Identifier, s)
}

func errorState(sm []string) *Token {
	err := fmt.Errorf(ErrUnexpectedCharacter, sm[0])
	return MakeToken(Error, data.String(err.Error()))
}
