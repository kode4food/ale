package lex

import (
	"fmt"
	"regexp"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/comb/slices"
)

type (
	Matcher func(src string) (*Token, string)

	tokenizer func([]string) *Token

	matchEntry struct {
		pattern  *regexp.Regexp
		function tokenizer
	}

	matchEntries []matchEntry
)

// Error messages
const (
	ErrStringNotTerminated  = "string has no closing quote"
	ErrUnexpectedCharacters = "unexpected characters: %s"
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

	Ignorable = matchEntries{
		pattern(lang.Comment, tokenState(Comment)),
		pattern(lang.NewLine, tokenState(NewLine)),
		pattern(lang.Whitespace, tokenState(Whitespace)),
	}

	Structure = matchEntries{
		pattern(lang.ListStart, tokenState(ListStart)),
		pattern(lang.VectorStart, tokenState(VectorStart)),
		pattern(lang.ObjectStart, tokenState(ObjectStart)),
		pattern(lang.ListEnd, tokenState(ListEnd)),
		pattern(lang.VectorEnd, tokenState(VectorEnd)),
		pattern(lang.ObjectEnd, tokenState(ObjectEnd)),
	}

	Quoting = matchEntries{
		pattern(lang.Quote, tokenState(QuoteMarker)),
		pattern(lang.SyntaxQuote, tokenState(SyntaxMarker)),
		pattern(lang.Splice, tokenState(SpliceMarker)),
		pattern(lang.Unquote, tokenState(UnquoteMarker)),
	}

	Values = matchEntries{
		pattern(lang.String, stringState),
		pattern(lang.Ratio, ratioState),
		pattern(lang.Float, floatState),
		pattern(lang.Integer, integerState),
	}

	Symbols = matchEntries{
		pattern(lang.Keyword, tokenState(Keyword)),
		pattern(lang.Identifier, identifierState),
	}
)

// Match tokenizes a String using the provided Matcher
func Match(src data.String, match Matcher) data.Sequence {
	var resolver sequence.LazyResolver
	var line, column int
	input := string(src)

	resolver = func() (data.Value, data.Sequence, bool) {
		if t, rest := match(input); t.Type() != endOfFile {
			t := t.withLocation(line, column)

			if t.isNewLine() {
				line++
				column = 0
			} else {
				column += len(input) - len(rest)
			}

			input = rest
			return t, sequence.NewLazy(resolver), true
		}
		return data.Null, data.Null, false
	}

	return sequence.NewLazy(resolver)
}

// StripWhitespace filters away all whitespace LangTokens from a Lexer Sequence
func StripWhitespace(s data.Sequence) data.Sequence {
	return sequence.Filter(s, func(v data.Value) bool {
		return !v.(*Token).isWhitespace()
	})
}

func MakeMatcher(m ...matchEntries) Matcher {
	entries := makeMatchEntries(m...)

	return func(src string) (*Token, string) {
		for _, s := range entries {
			if sm := s.pattern.FindStringSubmatch(src); sm != nil {
				return s.function(sm), src[len(sm[0]):]
			}
		}
		// Programmer error
		panic("unmatched lexing state")
	}
}

func makeMatchEntries(m ...matchEntries) matchEntries {
	var entries matchEntries
	entries = append(entries, pattern(`$`, endState(endOfFile)))
	for _, e := range m {
		entries = append(entries, e...)
	}
	entries = append(entries, pattern(lang.AnyChar, errorState))
	return entries
}

func (m matchEntries) Error() matchEntries {
	return slices.Map(m, func(e matchEntry) matchEntry {
		return matchEntry{
			pattern:  e.pattern,
			function: errorState,
		}
	})
}

func pattern(p string, s tokenizer) matchEntry {
	return matchEntry{
		pattern:  regexp.MustCompile("^" + p),
		function: s,
	}
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
	err := fmt.Errorf(ErrUnexpectedCharacters, sm[0])
	return MakeToken(Error, data.String(err.Error()))
}
