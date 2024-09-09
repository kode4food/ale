package data

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/kode4food/ale/internal/types"
)

// String is the Sequence-compatible representation of string values
type String string

// EmptyString represents the... empty string
const EmptyString = String("")

const (
	// ErrInvalidStartIndex is raised when an attempt to take a substring of a
	// String would result in a start index outside the String's actual range
	ErrInvalidStartIndex = "starting index is out of range: %d"

	// ErrInvalidEndIndex is raised when an attempt to take a substring of a
	// String would result in an end index beyond the String's actual range
	ErrInvalidEndIndex = "ending index is out of range: %d"

	// ErrEndIndexTooLow is raised when an attempt to take a substring of a
	// String specifies an ending index that is lower than the starting index
	ErrEndIndexTooLow = "ending index is lower than starting index: %d < %d"
)

var unescapeTable = map[string]string{
	"\\": "\\\\",
	"\n": "\\n",
	"\t": "\\t",
	"\f": "\\f",
	"\b": "\\b",
	"\r": "\\r",
	"\"": "\\\"",
}

var (
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
func (s String) Count() Integer {
	return Integer(utf8.RuneCountInString(string(s)))
}

// ElementAt returns the Character at the indexed position in the String
func (s String) ElementAt(index Integer) (Value, bool) {
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
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return String(runes[n:])
}

func (s String) Call(args ...Value) Value {
	idx := args[0].(Integer)
	if len(args) == 1 {
		return s.Call(idx, idx+1)
	} else if idx < 0 {
		panic(fmt.Errorf(ErrInvalidStartIndex, idx))
	}

	ns, ok := s.from(idx)
	if !ok {
		panic(fmt.Errorf(ErrInvalidStartIndex, idx))
	}

	end := args[1].(Integer)
	if end < idx {
		panic(fmt.Errorf(ErrEndIndexTooLow, idx, end))
	}
	if res, ok := ns.take(end - idx); ok {
		return res
	}
	panic(fmt.Errorf(ErrInvalidEndIndex, end))
}

func (s String) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (s String) from(index Integer) (String, bool) {
	ns := string(s)
	for i := Integer(0); i < index; i++ {
		_, w := utf8.DecodeRuneInString(ns)
		if w == 0 {
			return EmptyString, false
		}
		ns = ns[w:]
	}
	return String(ns), len(ns) > 0
}

func (s String) take(count Integer) (String, bool) {
	ns := string(s)
	for i := Integer(0); i < count; i++ {
		_, w := utf8.DecodeRuneInString(ns)
		if w == 0 {
			return EmptyString, false
		}
		ns = ns[w:]
	}
	return s[0 : len(s)-len(ns)], true
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
	buf.WriteString(`"`)
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		ch := string(f.(String))
		if res, ok := unescapeTable[ch]; ok {
			buf.WriteString(res)
		} else {
			buf.WriteString(ch)
		}
	}
	buf.WriteString(`"`)
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
