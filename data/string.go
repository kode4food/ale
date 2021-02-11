package data

import (
	"bytes"
	"unicode/utf8"
)

// String is the Sequence-compatible representation of string values
type String string

// EmptyString represents the... empty string
const EmptyString = String("")

var unescapeTable = map[string]string{
	"\\": "\\\\",
	"\n": "\\n",
	"\t": "\\t",
	"\f": "\\f",
	"\b": "\\b",
	"\r": "\\r",
	"\"": "\\\"",
}

// First returns the first character of the String
func (s String) First() Value {
	if r, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return String(r)
	}
	return Nil
}

// Rest returns a String of all characters after the first
func (s String) Rest() Sequence {
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
	return Nil, Nil, false
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
		return Nil, false
	}
	ns := string(s)
	p := 0
	for i := 0; i < index; i++ {
		if _, w := utf8.DecodeRuneInString(ns[p:]); w > 0 {
			p += w
		} else {
			return Nil, false
		}
	}
	if r, w := utf8.DecodeRuneInString(ns[p:]); w > 0 {
		return String(r), true
	}
	return Nil, false
}

// Equal compares this String to another for equality
func (s String) Equal(v Value) bool {
	if v, ok := v.(String); ok {
		return s == v
	}
	return false
}

// String converts this Value into a string
func (s String) String() string {
	return string(s)
}

// HashCode returns a hash code for the String
func (s String) HashCode() uint64 {
	return HashString(string(s))
}

// Quote quotes and escapes a string
func (s String) Quote() string {
	var buf bytes.Buffer
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

// MaybeQuoteString converts Values to strings, quoting wrapped Strings
func MaybeQuoteString(v Value) string {
	if s, ok := v.(String); ok {
		return s.Quote()
	}
	return v.String()
}
