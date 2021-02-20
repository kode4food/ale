package data

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

type (
	// Symbol is an identifier that can be resolved
	Symbol interface {
		symbol() // marker
		Value
		Named
	}

	// LocalSymbol represents an unqualified symbol that requires resolution
	LocalSymbol interface {
		localSymbol() // marker
		Symbol
	}

	// QualifiedSymbol represents a domain-qualified symbol
	QualifiedSymbol interface {
		Symbol
		Domain() Name
		Qualified() Name
	}

	// SymbolGenerator produces instance-unique local symbols
	SymbolGenerator struct {
		sync.Mutex
		data   [128]uint8
		maxPos int
	}

	localSymbol Name

	qualifiedSymbol struct {
		domain Name
		name   Name
	}
)

const (
	// DomainSeparator is the character used to separate a domain from
	// the local component of a qualified symbol
	DomainSeparator = '/'

	decimal        = "0123456789"
	lower          = "abcdefghijklmnopqrstuvwxyz"
	upper          = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base64Digits   = decimal + lower + upper + "%#"
	genSymTemplate = "x-%s-gensym-%s"
	genSymOverflow = uint8(len(base64Digits))
)

var gen = NewSymbolGenerator()

// NewGeneratedSymbol creates a generated Symbol
func NewGeneratedSymbol(name Name) Symbol {
	return gen.Local(name)
}

// ParseSymbol parses a qualified Name and produces a Symbol
func ParseSymbol(s String) Symbol {
	n := string(s)
	if i := strings.IndexRune(n, DomainSeparator); i > 0 {
		name := Name(n[i+1:])
		domain := Name(n[:i])
		return NewQualifiedSymbol(name, domain)
	}
	return localSymbol(s)
}

// NewSymbolGenerator creates a new symbol generator. In general, it is safe
// to use the global generator because it only maintains an incrementer
func NewSymbolGenerator() *SymbolGenerator {
	return new(SymbolGenerator)
}

// Local returns a newly generated local symbol
func (g *SymbolGenerator) Local(name Name) LocalSymbol {
	g.Lock()
	defer g.Unlock()
	g.inc(0)
	idx := g.str()
	q := fmt.Sprintf(genSymTemplate, name, idx)
	return localSymbol(q)
}

func (g *SymbolGenerator) inc(pos int) {
	if val := g.data[pos] + 1; val == genSymOverflow {
		g.overflow(pos)
	} else {
		g.data[pos] = val
	}
}

func (g *SymbolGenerator) overflow(pos int) {
	g.data[pos] = 0
	next := pos + 1
	if next > g.maxPos {
		g.maxPos = next
	}
	g.inc(next)
}

func (g *SymbolGenerator) str() string {
	var buf bytes.Buffer
	data := g.data
	for i := g.maxPos; i >= 0; i-- {
		d := base64Digits[data[i]]
		buf.WriteByte(d)
	}
	return buf.String()
}

// NewLocalSymbol returns a local symbol
func NewLocalSymbol(name Name) Symbol {
	return localSymbol(name)
}

func (localSymbol) symbol()      {}
func (localSymbol) localSymbol() {}

func (l localSymbol) Name() Name {
	return Name(l)
}

func (l localSymbol) Equal(v Value) bool {
	if v, ok := v.(localSymbol); ok {
		return l == v
	}
	return false
}

func (l localSymbol) String() string {
	return string(l)
}

func (l localSymbol) HashCode() uint64 {
	return HashString(string(l))
}

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name Name, domain Name) Symbol {
	return qualifiedSymbol{
		domain: domain,
		name:   name,
	}
}

func (qualifiedSymbol) symbol() {}

func (s qualifiedSymbol) Name() Name {
	return s.name
}

func (s qualifiedSymbol) Domain() Name {
	return s.domain
}

func (s qualifiedSymbol) Qualified() Name {
	var buf bytes.Buffer
	buf.WriteString(string(s.domain))
	buf.WriteRune(DomainSeparator)
	buf.WriteString(string(s.name))
	return Name(buf.String())
}

func (s qualifiedSymbol) Equal(v Value) bool {
	if v, ok := v.(qualifiedSymbol); ok {
		return s == v
	}
	return false
}

func (s qualifiedSymbol) String() string {
	return string(s.Qualified())
}

func (s qualifiedSymbol) HashCode() uint64 {
	return HashString(string(s.name)) * HashString(string(s.domain))
}
