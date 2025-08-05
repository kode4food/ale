package vm

import (
	"math/rand/v2"
	"slices"
	"sync/atomic"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/types"
)

// Procedure encapsulates the initial environment of an abstract machine
type Procedure struct {
	ArityChecker data.ArityChecker
	isa.Runnable
	hash atomic.Uint64
}

var (
	procSalt = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		data.Hashed
		data.Mapped
		data.Procedure
	} = (*Procedure)(nil)
)

func MakeProcedure(run *isa.Runnable, arity data.ArityChecker) *Procedure {
	return &Procedure{
		Runnable:     *run,
		ArityChecker: arity,
	}
}

// Call allows an abstract machine Procedure to be called to instantiate a
// Closure. Only the compiler invokes this calling interface.
func (p *Procedure) Call(values ...ale.Value) ale.Value {
	return &Closure{
		Procedure: p,
		captured:  slices.Clone(values),
	}
}

// CheckArity performs a compile-time arity check for the Procedure
func (p *Procedure) CheckArity(int) error {
	return nil
}

// Type makes Procedure a typed value
func (p *Procedure) Type() ale.Type {
	return types.MakeLiteral(types.BasicProcedure, p)
}

func (p *Procedure) Get(key ale.Value) (ale.Value, bool) {
	return data.DumpMapped(p).Get(key)
}

// Equal compares this Procedure to another for equality
func (p *Procedure) Equal(other ale.Value) bool {
	if other, ok := other.(*Procedure); ok {
		if p == other {
			return true
		}
		if p.StackSize != other.StackSize ||
			p.LocalCount != other.LocalCount ||
			p.Globals != other.Globals {
			return false
		}
		return basics.Equal(p.Code, other.Code) &&
			p.Constants.Equal(other.Constants)
	}
	return false
}

func (p *Procedure) HashCode() uint64 {
	if h := p.hash.Load(); h != 0 {
		return h
	}
	res := procSalt
	for i, inst := range p.Code {
		res ^= uint64(inst + 1)
		res ^= data.HashInt(i)
	}
	for i, c := range p.Constants {
		res ^= data.HashCode(c)
		res ^= data.HashInt(i)
	}
	p.hash.Store(res)
	return res
}
