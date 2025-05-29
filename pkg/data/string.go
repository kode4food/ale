package data

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// String is the Sequence-compatible representation of string data
type String string

// EmptyString represents the... empty string
const EmptyString = String("")

const (
	// ErrInvalidIndexes is raised when an attempt to take a substring receives
	// invalid start and end indexes
	ErrInvalidIndexes = "%d and %d are not valid start/end indices"

	// ErrInvalidStartIndex is raised when an attempt to take a substring of a
	// String specifies an ending index that is lower than the starting index
	ErrInvalidStartIndex = "%d is an invalid start index"
)

var (
	unescapeTable = map[string]string{
		lang.StringQuote: `\` + lang.StringQuote,

		"\b": `\b`,
		"\f": `\f`,
		"\n": `\n`,
		"\r": `\r`,
		"\t": `\t`,
		"\\": `\\`,
	}

	// compile-time checks for interface implementation
	_ interface {
		Caller
		Hashed
		RandomAccess
		Reverser
		Typed
	} = EmptyString
)

// Car returns the first character of the String
func (s String) Car() Value {
	if r, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return String(r)
	}
	return Null
}

// Cdr returns a String of all characters after the first
func (s String) Cdr() Value {
	if _, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return s[w:]
	}
	return EmptyString
}

// Split returns the split form (First and Rest) of the Sequence
func (s String) Split() (Value, Sequence, bool) {
	if r, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return String(r), s[w:], true
	}
	return Null, Null, false
}

// IsEmpty returns whether this sequence is empty
func (s String) IsEmpty() bool {
	return len(s) == 0
}

// Count returns the length of the String
func (s String) Count() int {
	return utf8.RuneCountInString(string(s))
}

// ElementAt returns the Character at the indexed position in the String
func (s String) ElementAt(index int) (Value, bool) {
	if index < 0 {
		return Null, false
	}
	ns, ok := s.from(index)
	if !ok {
		return Null, false
	}
	if r, w := utf8.DecodeRuneInString(string(ns)); w > 0 {
		return String(r), true
	}
	return Null, false
}

func (s String) Reverse() Sequence {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= 1 {
		return s
	}
	slices.Reverse(runes)
	return String(runes)
}

func (s String) Call(args ...Value) Value {
	if len(args) == 1 {
		return s.callFrom(int(args[0].(Integer)))
	}
	return s.callRange(int(args[0].(Integer)), int(args[1].(Integer)))
}

func (s String) callFrom(idx int) Value {
	if idx < 0 {
		panic(fmt.Errorf(ErrInvalidStartIndex, idx))
	}
	if ns, ok := s.from(idx); ok {
		return ns
	}
	panic(fmt.Errorf(ErrInvalidStartIndex, idx))
}

func (s String) callRange(idx, end int) Value {
	if idx < 0 || end < idx {
		panic(fmt.Errorf(ErrInvalidIndexes, idx, end))
	}

	ns, ok := s.from(idx)
	if !ok || len(ns) == 0 && end > idx {
		panic(fmt.Errorf(ErrInvalidIndexes, idx, end))
	}

	if res, ok := ns.take(end - idx); ok {
		return res
	}
	panic(fmt.Errorf(ErrInvalidIndexes, idx, end))
}

func (s String) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (s String) from(idx int) (String, bool) {
	_, r, ok := s.splitAt(idx)
	return r, ok
}

func (s String) take(count int) (String, bool) {
	f, _, ok := s.splitAt(count)
	return f, ok
}

func (s String) splitAt(idx int) (String, String, bool) {
	if idx == 0 {
		return EmptyString, s, true
	}
	ns := string(s)
	var r int
	for range idx {
		if _, w := utf8.DecodeRuneInString(ns[r:]); w > 0 {
			r += w
			continue
		}
		return EmptyString, EmptyString, false
	}
	return s[:r], s[r:], true
}

// Equal compares this String to another for equality
func (s String) Equal(other Value) bool {
	return s == other
}

// String converts this Value into a string
func (s String) String() string {
	return string(s)
}

// Type returns the Type for this String Value
func (String) Type() types.Type {
	return types.BasicString
}

// HashCode returns a hash code for the String
func (s String) HashCode() uint64 {
	return HashString(string(s))
}

// Quote quotes and escapes a string
func (s String) Quote() string {
	var buf strings.Builder
	buf.WriteString(lang.StringQuote)
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		ch := string(f.(String))
		if res, ok := unescapeTable[ch]; ok {
			buf.WriteString(res)
		} else {
			buf.WriteString(ch)
		}
	}
	buf.WriteString(lang.StringQuote)
	return buf.String()
}

func ToString(v Value) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return DumpString(v)
}

// ToQuotedString converts Values to string, possibly quoting wrapped Strings
func ToQuotedString(v Value) string {
	if s, ok := v.(String); ok {
		return s.Quote()
	}
	return ToString(v)
}
