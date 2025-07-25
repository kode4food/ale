package data

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

type (
	// Symbol is an identifier that can be resolved
	Symbol interface {
		ale.Typed
		Local() Local
	}

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

	// Local represents an unqualified symbol that requires resolution
	Local Name

	// Locals represent a set of Local
	Locals []Local

	qualified struct {
		domain Local
		name   Local
	}
)

const (
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

	lclSalt  = rand.Uint64()
	qualSalt = rand.Uint64()

	qualRegex = regexp.MustCompile("^" + lang.Qualified + "$")
	lclRegex  = regexp.MustCompile("^" + lang.Local + "$")

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Symbol
		fmt.Stringer
	} = Local("")

	_ interface {
		Hashed
		Qualified
		fmt.Stringer
	} = (*qualified)(nil)
)

// NewGeneratedSymbol creates a generated Symbol
func NewGeneratedSymbol(name Local) Symbol {
	return gen.Local(name)
}

// ParseSymbol parses a qualified Local and produces a Symbol
func ParseSymbol(s String) (Symbol, error) {
	n := string(s)
	if qualRegex.MatchString(n) {
		i := strings.Index(n, lang.DomainSeparator)
		name := Local(n[i+len(lang.DomainSeparator):])
		domain := Local(n[:i])
		res := NewQualifiedSymbol(name, domain)
		return res, nil
	}
	if lclRegex.MatchString(n) {
		return Local(s), nil
	}
	return nil, fmt.Errorf(ErrInvalidSymbol, n)
}

// MustParseSymbol parses a qualified Local and produces a Symbol or explodes
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
	defer g.Unlock()
	idx := g.str()
	g.inc(0)
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

func (l Local) Local() Local {
	return l
}

func (l Local) Equal(other ale.Value) bool {
	return l == other
}

func (Local) Type() ale.Type {
	return types.BasicSymbol
}

func (l Local) String() string {
	return string(l)
}

func (l Local) HashCode() uint64 {
	return lclSalt ^ HashString(string(l))
}

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name Local, domain Local) Symbol {
	return qualified{
		domain: domain,
		name:   name,
	}
}

func (s qualified) Local() Local {
	return s.name
}

func (s qualified) Domain() Local {
	return s.domain
}

func (s qualified) Equal(other ale.Value) bool {
	return s == other
}

func (qualified) Type() ale.Type {
	return types.BasicSymbol
}

func (s qualified) String() string {
	return string(s.domain) + lang.DomainSeparator + string(s.name)
}

func (s qualified) HashCode() uint64 {
	return qualSalt ^ HashString(string(s.name)) ^ HashString(string(s.domain))
}
