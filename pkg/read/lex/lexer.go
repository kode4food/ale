package lex

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/comb/basics"
)

type (
	Matcher func(string) (*Token, string)

	Matchers []Matcher

	tokenizer func([]string) *Token
)

const (
	// ErrStringNotTerminated is raised when the lexer reaches the end of its
	// stream while scanning a string, without encountering a closing quote
	ErrStringNotTerminated = "string has no closing quote"

	// ErrCommentNotTerminated is raised when the lexer reaches the end of its
	// stream while scanning a comment, without encountering an end marker
	ErrCommentNotTerminated = "comment has no closing comment marker"

	// ErrUnexpectedCharacters is raised when the lexer encounters a set of
	// characters that don't match any of the defined scanning patterns
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

	blockCommentStart = regexp.MustCompile(`^` + lang.BlockCommentStart)
	blockCommentEnd   = regexp.MustCompile(`^` + lang.BlockCommentEnd)

	Ignorable = Matchers{
		blockCommentMatcher,
		pattern(lang.Comment, tokenState(Comment)),
		pattern(lang.NewLine, tokenState(NewLine)),
		pattern(lang.Whitespace, tokenState(Whitespace)),
	}

	Structure = Matchers{
		pattern(lang.ListStart, tokenState(ListStart)),
		pattern(lang.VectorStart, tokenState(VectorStart)),
		pattern(lang.ObjectStart, tokenState(ObjectStart)),
		pattern(lang.ListEnd, tokenState(ListEnd)),
		pattern(lang.VectorEnd, tokenState(VectorEnd)),
		pattern(lang.ObjectEnd, tokenState(ObjectEnd)),
	}

	Quoting = Matchers{
		pattern(lang.Quote, tokenState(QuoteMarker)),
		pattern(lang.SyntaxQuote, tokenState(SyntaxMarker)),
		pattern(lang.Splice, tokenState(SpliceMarker)),
		pattern(lang.Unquote, tokenState(UnquoteMarker)),
	}

	Values = Matchers{
		pattern(lang.String, stringState),
		pattern(lang.Ratio, ratioState),
		pattern(lang.Float, floatState),
		pattern(lang.Integer, integerState),
	}

	Symbols = Matchers{
		pattern(lang.Keyword, tokenState(Keyword)),
		pattern(lang.Identifier, identifierState),
	}
)

// Match tokenizes a String using the provided Matcher
func Match(src data.String, m Matcher) data.Sequence {
	var resolver sequence.LazyResolver
	var line, column int
	input := string(src)

	resolver = func() (data.Value, data.Sequence, bool) {
		if t, rest := m(input); t.Type() != endOfFile {
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

func ExhaustiveMatcher(all ...Matchers) Matcher {
	matchers := makeExhaustive(all...)
	return func(input string) (*Token, string) {
		for _, m := range matchers {
			if t, rest := m(input); t != nil {
				return t, rest
			}
		}
		panic(debug.ProgrammerError("unmatched lexing state"))
	}
}

func makeExhaustive(all ...Matchers) Matchers {
	cat := slices.Concat(all...)
	res := make(Matchers, 0, len(cat)+2)
	res = append(res, pattern(`$`, endState(endOfFile)))
	res = append(res, cat...)
	res = append(res, pattern(lang.AnyChar, errorState))
	return res
}

func (m Matchers) Error() Matchers {
	return basics.Map(m, func(wrapped Matcher) Matcher {
		return func(input string) (*Token, string) {
			if t, rest := wrapped(input); t != nil {
				err := fmt.Errorf(ErrUnexpectedCharacters, t.input)
				return MakeToken(Error, data.String(err.Error())), rest
			}
			return nil, input
		}
	})
}

func pattern(p string, t tokenizer) Matcher {
	r := regexp.MustCompile("^" + p)
	return func(input string) (*Token, string) {
		if sm := r.FindStringSubmatch(input); sm != nil {
			return t(sm).withInput(sm[0]), input[len(sm[0]):]
		}
		return nil, input
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

func blockCommentMatcher(input string) (*Token, string) {
	start := blockCommentStart.FindStringSubmatch(input)
	if start == nil {
		return nil, input
	}

	stack := 1
	for i := len(start[0]); i < len(input); {
		next := input[i:]
		if sm := blockCommentStart.FindStringSubmatch(next); sm != nil {
			i += len(sm[0])
			stack++
			continue
		}
		if sm := blockCommentEnd.FindStringSubmatch(next); sm != nil {
			i += len(sm[0])
			stack--
			if stack == 0 {
				c := input[:i]
				rest := input[i:]
				return MakeToken(Comment, data.String(c)), rest
			}
			continue
		}
		i++
	}
	return MakeToken(Error, data.String(ErrCommentNotTerminated)), input
}
