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

	tokenizer func([]string) *Token
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
	ErrUnmatchedComment = "encountered '|#' with no open block comment"

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

	dotRegex          = regexp.MustCompile(`^` + lang.Dot + `$`)
	blockCommentStart = regexp.MustCompile(`^` + lang.BlockCommentStart)
	blockCommentEnd   = regexp.MustCompile(`^` + lang.BlockCommentEnd)
	newLine           = regexp.MustCompile(lang.NewLine)

	Ignorable = Matchers{
		blockCommentMatcher,
		patternMatcher(lang.Comment, tokenState(Comment)),
		patternMatcher(lang.NewLine, tokenState(NewLine)),
		patternMatcher(lang.Whitespace, tokenState(Whitespace)),
	}

	Structure = Matchers{
		patternMatcher(lang.ListStart, tokenState(ListStart)),
		patternMatcher(lang.VectorStart, tokenState(VectorStart)),
		patternMatcher(lang.ObjectStart, tokenState(ObjectStart)),
		patternMatcher(lang.ListEnd, tokenState(ListEnd)),
		patternMatcher(lang.VectorEnd, tokenState(VectorEnd)),
		patternMatcher(lang.ObjectEnd, tokenState(ObjectEnd)),
	}

	Quoting = Matchers{
		patternMatcher(lang.Quote, tokenState(QuoteMarker)),
		patternMatcher(lang.SyntaxQuote, tokenState(SyntaxMarker)),
		patternMatcher(lang.Splice, tokenState(SpliceMarker)),
		patternMatcher(lang.Unquote, tokenState(UnquoteMarker)),
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
	input := string(src)

	resolver = func() (data.Value, data.Sequence, bool) {
		if t, rest := m(input); t.Type() != endOfFile {
			t := t.withLocation(line, column)
			line, column = bumpLocation(t.input, line, column)
			input = rest
			return t, sequence.NewLazy(resolver), true
		}
		return data.Null, data.Null, false
	}

	return sequence.NewLazy(resolver)
}

// StripWhitespace filters away all whitespace Tokens from a Lexer Sequence
func StripWhitespace(s data.Sequence) data.Sequence {
	return sequence.Filter(s, func(v data.Value) bool {
		return !v.(*Token).isWhitespace()
	})
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
	res = append(res, patternMatcher(`$`, endState(endOfFile)))
	res = append(res, cat...)
	res = append(res, patternMatcher(lang.AnyChar, errorState))
	return res
}

func bumpLocation(i string, l, c int) (int, int) {
	s := newLine.Split(i, -1)
	dl := len(s) - 1
	dc := len(s[len(s)-1])
	if dl > 0 {
		return l + dl, dc
	}
	return l, c + dc
}

func patternMatcher(p string, t tokenizer) Matcher {
	r := regexp.MustCompile("^" + p)
	return func(input string) (*Token, string) {
		if sm := r.FindStringSubmatch(input); sm != nil {
			return t(sm).withInput(sm[0]), input[len(sm[0]):]
		}
		return nil, input
	}
}

func endState(t TokenType) tokenizer {
	return func([]string) *Token {
		return MakeToken(t, nil)
	}
}

func tokenState(t TokenType) tokenizer {
	return func(_ []string) *Token {
		return MakeToken(t, data.Null)
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

func keywordState(sm []string) *Token {
	s := strings.Clone(sm[0])
	return MakeToken(Keyword, data.String(s))
}

func identifierState(sm []string) *Token {
	s := strings.Clone(sm[0])
	if dotRegex.MatchString(s) {
		return MakeToken(Dot, data.String(s))
	}
	return MakeToken(Identifier, data.String(s))
}

func errorState(sm []string) *Token {
	s := strings.Clone(sm[0])
	err := fmt.Errorf(ErrUnexpectedCharacters, s)
	return MakeToken(Error, data.String(err.Error()))
}

func blockCommentMatcher(input string) (*Token, string) {
	start := blockCommentStart.FindStringSubmatch(input)
	if start == nil {
		return blockCommentStrayEndMatcher(input)
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
				t := MakeToken(BlockComment, data.Null)
				return t.withInput(input[:i]), input[i:]
			}
			continue
		}
		i++
	}
	return MakeToken(Error, data.String(ErrCommentNotTerminated)), input
}

func blockCommentStrayEndMatcher(input string) (*Token, string) {
	if sm := blockCommentEnd.FindStringSubmatch(input); sm != nil {
		err := data.String(ErrUnmatchedComment)
		rest := input[len(sm[0]):]
		return MakeToken(Error, err).withInput(sm[0]), rest
	}
	return nil, input
}
