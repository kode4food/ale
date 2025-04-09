package lex

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/comb/basics"
)

type (
	Matcher func(string) (*Token, string)

	Matchers []Matcher

	tokenizer func(string) *Token
)

const (
	// ErrStringNotTerminated is raised when the lexer reaches the end of its
	// stream while scanning a string, without encountering a closing quote
	ErrStringNotTerminated = "string has no closing quote"

	// ErrCommentNotTerminated is raised when the lexer reaches the end of its
	// stream while scanning a comment, without encountering an end marker
	ErrCommentNotTerminated = "comment has no closing comment marker"

	// ErrUnmatchedComment is raised when a comment end marker is encountered
	// in the stream when no open block comment is being parsed
	ErrUnmatchedComment = "encountered '" + lang.BlockCommentEnd +
		"' with no open block comment"

	// ErrUnexpectedCharacters is raised when the lexer encounters a set of
	// characters that don't match any of the defined scanning patterns
	ErrUnexpectedCharacters = "unexpected characters: %s"
)

var (
	escaped    = regexp.MustCompile(`\\\\|\\"|\\[^\\"]`)
	escapedMap = map[string]string{
		`\b`: "\b",
		`\f`: "\f",
		`\n`: "\n",
		`\r`: "\r",
		`\s`: " ",
		`\t`: "\t",
	}

	newLine = regexp.MustCompile(lang.NewLine)

	Ignorable = Matchers{
		blockCommentMatcher,
		patternMatcher(lang.Comment, tokenState(Comment)),
		patternMatcher(lang.NewLine, tokenState(NewLine)),
		patternMatcher(lang.Whitespace, tokenState(Whitespace)),
	}

	Structure = Matchers{
		prefixMatcher(lang.ListStart, tokenState(ListStart)),
		prefixMatcher(lang.VectorStart, tokenState(VectorStart)),
		prefixMatcher(lang.ObjectStart, tokenState(ObjectStart)),
		prefixMatcher(lang.ListEnd, tokenState(ListEnd)),
		prefixMatcher(lang.VectorEnd, tokenState(VectorEnd)),
		prefixMatcher(lang.ObjectEnd, tokenState(ObjectEnd)),
	}

	Quoting = Matchers{
		prefixMatcher(lang.Quote, tokenState(QuoteMarker)),
		prefixMatcher(lang.SyntaxQuote, tokenState(SyntaxMarker)),
		prefixMatcher(lang.Splice, tokenState(SpliceMarker)),
		prefixMatcher(lang.Unquote, tokenState(UnquoteMarker)),
	}

	Values = Matchers{
		patternMatcher(lang.String, stringState),
		patternMatcher(lang.Ratio, ratioState),
		patternMatcher(lang.Float, floatState),
		patternMatcher(lang.Integer, integerState),
	}

	Symbols = Matchers{
		patternMatcher(lang.Keyword, keywordState),
		patternMatcher(lang.Identifier, identifierState),
	}
)

// Match tokenizes a String using the provided Matcher
func Match(src data.String, m Matcher) data.Sequence {
	var resolver sequence.LazyResolver
	var line, column int
	var t *Token
	input := string(src)

	resolver = func() (data.Value, data.Sequence, bool) {
		t, input = m(input)
		t.SetLocation(line, column)
		if t.Type() != EOF {
			line, column = bumpLocation(t.input, line, column)
			return t, sequence.NewLazy(resolver), true
		}
		return data.Null, data.Null, false
	}

	return sequence.NewLazy(resolver)
}

// StripWhitespace filters away all whitespace Tokens from a Lexer Sequence
func StripWhitespace(s data.Sequence) data.Sequence {
	return sequence.Filter(s, func(v data.Value) bool {
		return !v.(*Token).IsWhitespace()
	})
}

func (m Matchers) Error() Matchers {
	return basics.Map(m, func(wrapped Matcher) Matcher {
		return func(input string) (*Token, string) {
			if t, rest := wrapped(input); t != nil {
				err := fmt.Errorf(ErrUnexpectedCharacters, t.input)
				s := data.String(err.Error())
				return Error.FromValue(t.input, s), rest
			}
			return nil, input
		}
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
	res = append(res, exactMatcher(``, tokenState(EOF)))
	res = append(res, cat...)
	res = append(res, anyCharMatcher(errorState))
	return res
}

func bumpLocation(i string, l, c int) (int, int) {
	s := newLine.Split(i, -1)
	dl := len(s) - 1
	dc := len(s[dl])
	if dl == 0 { // no newline, only bump column
		return l, c + dc
	}
	// bump line and reset column
	return l + dl, dc
}

func anyCharMatcher(t tokenizer) Matcher {
	return func(input string) (*Token, string) {
		if len(input) > 0 {
			return t(strings.Clone(input[:1])), input[1:]
		}
		return nil, input
	}
}

func exactMatcher(p string, t tokenizer) Matcher {
	token := t(p)
	return func(input string) (*Token, string) {
		if input == p {
			return token, ``
		}
		return nil, input
	}
}

func prefixMatcher(p string, t tokenizer) Matcher {
	token := t(p)
	return func(input string) (*Token, string) {
		if strings.HasPrefix(input, p) {
			return token, input[len(p):]
		}
		return nil, input
	}
}

func patternMatcher(p string, t tokenizer) Matcher {
	r := regexp.MustCompile("^" + p)
	return func(input string) (*Token, string) {
		if sm := r.FindStringSubmatch(input); sm != nil {
			return t(strings.Clone(sm[0])), input[len(sm[0]):]
		}
		return nil, input
	}
}

func tokenState(typ TokenType) tokenizer {
	return func(m string) *Token {
		return typ.From(m)
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

func stringState(m string) *Token {
	eos := len(m) - 1
	if len(m) <= 1 || m[eos] != '"' {
		err := data.String(ErrStringNotTerminated)
		return Error.FromValue(m, err)
	}
	s := unescape(m[1:eos])
	return String.FromValue(m, data.String(s))
}

func ratioState(m string) *Token {
	res, err := data.ParseRatio(m)
	return tokenizeNumber(m, res, err)
}

func floatState(m string) *Token {
	res, err := data.ParseFloat(m)
	return tokenizeNumber(m, res, err)
}

func integerState(m string) *Token {
	res, err := data.ParseInteger(m)
	return tokenizeNumber(m, res, err)
}

func tokenizeNumber(m string, res data.Number, err error) *Token {
	if err != nil {
		return Error.FromValue(m, data.String(err.Error()))
	}
	return Number.FromValue(m, res)
}

func keywordState(m string) *Token {
	return Keyword.FromValue(m, data.String(m))
}

func identifierState(m string) *Token {
	if m == lang.Dot {
		return Dot.From(m)
	}
	return Identifier.FromValue(m, data.String(m))
}

func errorState(m string) *Token {
	err := fmt.Errorf(ErrUnexpectedCharacters, m)
	return Error.FromValue(m, data.String(err.Error()))
}

func blockCommentMatcher(input string) (*Token, string) {
	if !strings.HasPrefix(input, lang.BlockCommentStart) {
		return blockCommentStrayEndMatcher(input)
	}

	stack := 1
	for i := len(lang.BlockCommentStart); i < len(input); {
		next := input[i:]
		if strings.HasPrefix(next, lang.BlockCommentStart) {
			i += len(lang.BlockCommentStart)
			stack++
			continue
		}
		if strings.HasPrefix(next, lang.BlockCommentEnd) {
			i += len(lang.BlockCommentEnd)
			stack--
			if stack == 0 {
				m := strings.Clone(input[:i])
				return BlockComment.From(m), input[i:]
			}
			continue
		}
		i++
	}
	err := data.String(ErrCommentNotTerminated)
	return Error.FromValue(input, err), input
}

func blockCommentStrayEndMatcher(input string) (*Token, string) {
	if strings.HasPrefix(input, lang.BlockCommentEnd) {
		err := data.String(ErrUnmatchedComment)
		t := Error.FromValue(lang.BlockCommentEnd, err)
		rest := input[len(lang.BlockCommentEnd):]
		return t, rest
	}
	return nil, input
}
