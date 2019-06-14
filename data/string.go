package data

import (
	"bytes"
	"unicode/utf8"
)

// String is the Sequence-compatible representation of string values
type String string

const emptyStr = String("")

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
	return Null
}

// Rest returns a String of all characters after the first
func (s String) Rest() Sequence {
	if _, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return String(s[w:])
	}
	return emptyStr
}

// Split returns the split form (First and Rest) of the Sequence
func (s String) Split() (Value, Sequence, bool) {
	if r, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return String(r), String(s[w:]), true
	}
	return Null, emptyStr, false
}

// IsEmpty returns whether or not this sequence is empty
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
	ns := string(s)
	p := 0
	for i := 0; i < index; i++ {
		if _, w := utf8.DecodeRuneInString(ns[p:]); w > 0 {
			p += w
		} else {
			return Null, false
		}
	}
	if r, w := utf8.DecodeRuneInString(ns[p:]); w > 0 {
		return String(r), true
	}
	return Null, false
}

// String converts this Value into a string
func (s String) String() string {
	return string(s)
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
