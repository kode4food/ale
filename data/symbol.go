package data

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/kode4food/ale/types"
)

type (
	// Symbol is an identifier that can be resolved
	Symbol interface {
		symbol() // marker
		Value
		Named
	}

	// LocalSymbols represents a set of LocalSymbol
	LocalSymbols []LocalSymbol

	// Named is the generic interface for values that are named
	Named interface {
		Name() LocalSymbol
	}

	// LocalSymbol represents an unqualified symbol that requires resolution
	LocalSymbol string

	// QualifiedSymbol represents a domain-qualified symbol
	QualifiedSymbol interface {
		Symbol
		Domain() LocalSymbol
	}

	// SymbolGenerator produces instance-unique local symbols
	SymbolGenerator struct {
		sync.Mutex
		data   [128]uint8
		maxPos int
	}

	qualifiedSymbol struct {
		domain LocalSymbol
		name   LocalSymbol
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
func NewGeneratedSymbol(name LocalSymbol) Symbol {
	return gen.Local(name)
}

// ParseSymbol parses a qualified Name and produces a Symbol
func ParseSymbol(s String) Symbol {
	n := string(s)
	if i := strings.IndexRune(n, DomainSeparator); i > 0 {
		name := LocalSymbol(n[i+1:])
		domain := LocalSymbol(n[:i])
		return NewQualifiedSymbol(name, domain)
	}
	return LocalSymbol(s)
}

// NewSymbolGenerator creates a new symbol generator. In general, it is safe to
// use the global generator because it only maintains an incrementer
func NewSymbolGenerator() *SymbolGenerator {
	return new(SymbolGenerator)
}

// Local returns a newly generated local symbol
func (g *SymbolGenerator) Local(name LocalSymbol) LocalSymbol {
	g.Lock()
	defer g.Unlock()
	g.inc(0)
	idx := g.str()
	q := fmt.Sprintf(genSymTemplate, name, idx)
	return LocalSymbol(q)
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

func (LocalSymbol) symbol() {}

func (l LocalSymbol) Name() LocalSymbol {
	return l
}

func (l LocalSymbol) Equal(v Value) bool {
	if v, ok := v.(LocalSymbol); ok {
		return l == v
	}
	return false
}

func (LocalSymbol) Type() types.Type {
	return types.Symbol
}

func (l LocalSymbol) String() string {
	return string(l)
}

func (l LocalSymbol) HashCode() uint64 {
	return HashString(string(l))
}

// Sorted returns a sorted set of LocalSymbols
func (n LocalSymbols) Sorted() LocalSymbols {
	res := make(LocalSymbols, len(n))
	copy(res, n)
	slices.Sort(res)
	return res
}

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name LocalSymbol, domain LocalSymbol) Symbol {
	return qualifiedSymbol{
		domain: domain,
		name:   name,
	}
}

func (qualifiedSymbol) symbol() {}

func (s qualifiedSymbol) Name() LocalSymbol {
	return s.name
}

func (s qualifiedSymbol) Domain() LocalSymbol {
	return s.domain
}

func (s qualifiedSymbol) Equal(v Value) bool {
	if v, ok := v.(qualifiedSymbol); ok {
		return s == v
	}
	return false
}

func (qualifiedSymbol) Type() types.Type {
	return types.Symbol
}

func (s qualifiedSymbol) String() string {
	var buf bytes.Buffer
	buf.WriteString(string(s.domain))
	buf.WriteRune(DomainSeparator)
	buf.WriteString(string(s.name))
	return buf.String()
}

func (s qualifiedSymbol) HashCode() uint64 {
	return HashString(string(s.name)) * HashString(string(s.domain))
}
