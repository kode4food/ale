package vm

import (
	"math/rand"
	"slices"
	"sync/atomic"

	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

// Procedure encapsulates the initial environment of an abstract machine
type Procedure struct {
	ArityChecker data.ArityChecker
	isa.Runnable
	hash atomic.Uint64
}

var procedureHash = rand.Uint64()

func MakeProcedure(run *isa.Runnable, arity data.ArityChecker) *Procedure {
	return &Procedure{
		Runnable:     *run,
		ArityChecker: arity,
	}
}

// Call allows an abstract machine Procedure to be called to instantiate a
// Closure. Only the compiler invokes this calling interface.
func (p *Procedure) Call(values ...data.Value) data.Value {
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
func (p *Procedure) Type() types.Type {
	return types.BasicProcedure
}

func (p *Procedure) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(p).Get(key)
}

// Equal compares this Procedure to another for equality
func (p *Procedure) Equal(other data.Value) bool {
	if other, ok := other.(*Procedure); ok {
		return p == other ||
			p.Globals == other.Globals &&
				p.StackSize == other.StackSize &&
				p.LocalCount == other.LocalCount &&
				slices.Equal(p.Code, other.Code) &&
				p.Constants.Equal(other.Constants)
	}
	return false
}

func (p *Procedure) HashCode() uint64 {
	if h := p.hash.Load(); h != 0 {
		return h
	}
	res := procedureHash
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
