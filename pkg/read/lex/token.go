package lex

import (
	"fmt"

	"github.com/kode4food/ale/pkg/data"
)

type (
	// TokenType is an opaque type for lexer Match
	TokenType int

	// Token is a lexer value
	Token struct {
		input  string
		value  data.Value
		typ    TokenType
		line   int
		column int
	}
)

const (
	Error TokenType = iota
	Keyword
	Identifier
	Dot
	String
	Number
	ListStart
	ListEnd
	VectorStart
	VectorEnd
	ObjectStart
	ObjectEnd
	QuoteMarker
	SyntaxMarker
	UnquoteMarker
	SpliceMarker
	Whitespace
	NewLine
	Comment
	endOfFile
)

const errTokenWrapped = "%w (line %d, column %d)"

// MakeToken constructs a new scanner Token
func MakeToken(t TokenType, v data.Value) *Token {
	return &Token{
		typ:   t,
		value: v,
	}
}

// withInput returns a copy of the Token with the input string
func (t *Token) withInput(i string) *Token {
	res := *t
	res.input = i
	return &res
}

// withLocation returns a copy of the Token with location information
func (t *Token) withLocation(line, column int) *Token {
	res := *t
	res.line = line
	res.column = column
	return &res
}

// Type returns the TokenType for this Token
func (t *Token) Type() TokenType {
	return t.typ
}

// Value returns the scanned data.Value for this Token
func (t *Token) Value() data.Value {
	return t.value
}

// Line returns the line where this Token occurs
func (t *Token) Line() int {
	return t.line
}

// Column returns the column where this Token occurs
func (t *Token) Column() int {
	return t.column
}

// Equal compares this Token to another for equality
func (t *Token) Equal(other data.Value) bool {
	if other, ok := other.(*Token); ok {
		return t == other || t.typ == other.typ && t.value.Equal(other.value)
	}
	return false
}

func (t *Token) isNewLine() bool {
	return t.typ == Comment || t.typ == NewLine
}

func (t *Token) isWhitespace() bool {
	return t.typ == Comment || t.typ == NewLine || t.typ == Whitespace
}

func (t *Token) WrapError(e error) error {
	return fmt.Errorf(errTokenWrapped, e, t.line+1, t.column+1)
}
