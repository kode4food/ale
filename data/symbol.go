package data

import (
	"bytes"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/types"
)

type (
	// Symbol is an identifier that can be resolved
	Symbol interface {
		symbol() // marker
		Value
		Named
	}

	// Locals represents a set of Local
	Locals []Local

	// Named is the generic interface for values that are named
	Named interface {
		Name() Local
	}

	// Local represents an unqualified symbol that requires resolution
	Local string

	// Qualified represents a domain-qualified symbol
	Qualified interface {
		Symbol
		Domain() Local
	}

	// SymbolGenerator produces instance-unique local symbols
	SymbolGenerator struct {
		sync.Mutex
		data   [128]uint8
		maxPos int
	}

	qualifiedSymbol struct {
		domain Local
		name   Local
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

// Error messages
const (
	ErrInvalidSymbol          = "invalid symbol: %s"
	ErrInvalidQualifiedSymbol = "invalid qualified symbol: %s"
)

var (
	gen = NewSymbolGenerator()

	symbolRegex = regexp.MustCompile("^" + lang.ID + "$")
)

// NewGeneratedSymbol creates a generated Symbol
func NewGeneratedSymbol(name Local) Symbol {
	return gen.Local(name)
}

// ParseSymbol parses a qualified Name and produces a Symbol
func ParseSymbol(s String) (Symbol, error) {
	n := string(s)
	if !symbolRegex.MatchString(n) {
		return nil, fmt.Errorf(ErrInvalidSymbol, n)
	}
	if i := strings.IndexRune(n, DomainSeparator); i > 0 {
		name := Local(n[i+1:])
		domain := Local(n[:i])
		if len(domain) == 0 || len(name) == 0 {
			return nil, fmt.Errorf(ErrInvalidQualifiedSymbol, n)
		}
		return NewQualifiedSymbol(name, domain), nil
	}
	return Local(s), nil
}

// MustParseSymbol parses a qualified Name and produces a Symbol or explodes
func MustParseSymbol(s String) Symbol {
	sym, err := ParseSymbol(s)
	if err != nil {
		panic(err)
	}
	return sym
}

// NewSymbolGenerator creates a new symbol generator. In general, it is safe to
// use the global generator because it only maintains an incrementer
func NewSymbolGenerator() *SymbolGenerator {
	return new(SymbolGenerator)
}

// Local returns a newly generated local symbol
func (g *SymbolGenerator) Local(name Local) Local {
	g.Lock()
	defer g.Unlock()
	g.inc(0)
	idx := g.str()
	q := fmt.Sprintf(genSymTemplate, name, idx)
	return Local(q)
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

func (Local) symbol() {}

func (l Local) Name() Local {
	return l
}

func (l Local) Equal(v Value) bool {
	if v, ok := v.(Local); ok {
		return l == v
	}
	return false
}

func (Local) Type() types.Type {
	return types.Symbol
}

func (l Local) String() string {
	return string(l)
}

func (l Local) HashCode() uint64 {
	return HashString(string(l))
}

// Sorted returns a sorted set of Locals
func (n Locals) Sorted() Locals {
	res := make(Locals, len(n))
	copy(res, n)
	slices.Sort(res)
	return res
}

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name Local, domain Local) Symbol {
	return qualifiedSymbol{
		domain: domain,
		name:   name,
	}
}

func (qualifiedSymbol) symbol() {}

func (s qualifiedSymbol) Name() Local {
	return s.name
}

func (s qualifiedSymbol) Domain() Local {
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
