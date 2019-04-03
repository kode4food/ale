package api

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type (
	// Symbol is an identifier that can be resolved
	Symbol interface {
		Value
		Named
		Symbol()
	}

	// LocalSymbol represents an unqualified symbol that requires resolution
	LocalSymbol interface {
		Symbol
		LocalSymbol()
	}

	// QualifiedSymbol represents a domain-qualified symbol
	QualifiedSymbol interface {
		Symbol
		Domain() Name
		Qualified() Name
	}

	localSymbol Name

	qualifiedSymbol struct {
		domain Name
		name   Name
	}
)

const genSymTemplate = "x-%s-gensym-%x"

var genSymIncrement uint64

// NewGeneratedSymbol creates a generated Symbol
func NewGeneratedSymbol(name Name) Symbol {
	idx := atomic.AddUint64(&genSymIncrement, 1)
	q := fmt.Sprintf(genSymTemplate, name, idx)
	return localSymbol(q)
}

// ParseSymbol parses a qualified Name and produces a Symbol
func ParseSymbol(n Name) Symbol {
	if i := strings.IndexRune(string(n), '/'); i > 0 {
		return NewQualifiedSymbol(n[i+1:], n[:i])
	}
	return localSymbol(n)
}

// NewLocalSymbol returns a local symbol
func NewLocalSymbol(name Name) Symbol {
	return localSymbol(name)
}

func (localSymbol) Symbol()      {}
func (localSymbol) LocalSymbol() {}

func (l localSymbol) Name() Name {
	return Name(l)
}

func (l localSymbol) String() string {
	return string(l)
}

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name Name, domain Name) Symbol {
	return qualifiedSymbol{
		domain: domain,
		name:   name,
	}
}

func (qualifiedSymbol) Symbol() {}

func (s qualifiedSymbol) Name() Name {
	return s.name
}

func (s qualifiedSymbol) Domain() Name {
	return s.domain
}

func (s qualifiedSymbol) Qualified() Name {
	return Name(s.domain + "/" + s.name)
}

func (s qualifiedSymbol) String() string {
	return string(s.Qualified())
}
