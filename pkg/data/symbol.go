package data

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/comb/basics"
)

type (
	// Symbol is an identifier that can be resolved
	Symbol interface {
		symbol() // marker
		Typed
		Value
		Named
	}

	// Locals represent a set of Local
	Locals []Local

	// Named is the generic interface for data that are named
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
		prefix string
		data   [128]uint8
		sync.Mutex
		maxPos int
	}

	qualified struct {
		domain Local
		name   Local
	}
)

const (
	// DomainSeparator is the character used to separate a domain from
	// the local component of a qualified symbol
	DomainSeparator = '/'

	decimal         = "0123456789"
	lower           = "abcdefghijklmnopqrstuvwxyz"
	upper           = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SymbolGenDigits = decimal + lower + upper + "-+"

	genSymTemplate = "x-%s-gensym-%s-%s"
	genSymOverflow = uint8(len(SymbolGenDigits))
)

// ErrInvalidSymbol is raised when a call to ParseSymbol can't interpret its
// input as a proper Symbol name (local or qualified)
const ErrInvalidSymbol = "invalid symbol: %s"

var (
	gen = NewSymbolGenerator()

	qualifiedRegex = regexp.MustCompile("^" + lang.Qualified + "$")
	localRegex     = regexp.MustCompile("^" + lang.Local + "$")
)

// NewGeneratedSymbol creates a generated Symbol
func NewGeneratedSymbol(name Local) Symbol {
	return gen.Local(name)
}

// ParseSymbol parses a qualified Name and produces a Symbol
func ParseSymbol(s String) (Symbol, error) {
	n := string(s)
	if qualifiedRegex.MatchString(n) {
		i := strings.IndexRune(n, DomainSeparator)
		name := Local(n[i+1:])
		domain := Local(n[:i])
		res := NewQualifiedSymbol(name, domain)
		return res, nil
	}
	if localRegex.MatchString(n) {
		return Local(s), nil
	}
	return nil, fmt.Errorf(ErrInvalidSymbol, n)
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
	return &SymbolGenerator{
		prefix: strconv.FormatUint(uint64(rand.Uint32()), 36),
	}
}

func (g *SymbolGenerator) Prefix() string {
	return g.prefix
}

// Local returns a newly generated local symbol
func (g *SymbolGenerator) Local(name Local) Local {
	g.Lock()
	idx := g.str()
	g.inc(0)
	g.Unlock()
	q := fmt.Sprintf(genSymTemplate, name, g.prefix, idx)
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
	data, maxPos := g.data, g.maxPos
	res := make([]byte, maxPos+1)
	for i := 0; i <= maxPos; i++ {
		res[i] = SymbolGenDigits[data[maxPos-i]]
	}
	return string(res)
}

func (Local) symbol() {}

func (l Local) Name() Local {
	return l
}

func (l Local) Equal(other Value) bool {
	return l == other
}

func (Local) Type() types.Type {
	return types.BasicSymbol
}

func (l Local) String() string {
	return string(l)
}

func (l Local) HashCode() uint64 {
	return HashString(string(l))
}

// Sorted returns a sorted set of Locals
func (n Locals) Sorted() Locals {
	return basics.Sort(n)
}

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name Local, domain Local) Symbol {
	return qualified{
		domain: domain,
		name:   name,
	}
}

func (qualified) symbol() {}

func (s qualified) Name() Local {
	return s.name
}

func (s qualified) Domain() Local {
	return s.domain
}

func (s qualified) Equal(other Value) bool {
	return s == other
}

func (qualified) Type() types.Type {
	return types.BasicSymbol
}

func (s qualified) String() string {
	var buf strings.Builder
	buf.WriteString(string(s.domain))
	buf.WriteRune(DomainSeparator)
	buf.WriteString(string(s.name))
	return buf.String()
}

func (s qualified) HashCode() uint64 {
	return HashString(string(s.name)) ^ HashString(string(s.domain))
}
