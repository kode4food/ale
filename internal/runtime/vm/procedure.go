package vm

import (
	"math/rand"
	"slices"
	"sync/atomic"

	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

type (
	// Procedure encapsulates the initial environment of an abstract machine
	Procedure struct {
		ArityChecker data.ArityChecker
		isa.Runnable
		hash  uint64
		flags ProcFlag
	}

	ProcFlag uint

	Closure struct {
		*Procedure
		captured data.Vector
		hash     uint64
	}
)

const (
	NoInline ProcFlag = 1 << iota
)

var procedureHash = rand.Uint64()

func MakeProcedure(run *isa.Runnable, arity data.ArityChecker) *Procedure {
	res := &Procedure{
		Runnable:     *run,
		ArityChecker: arity,
	}

	if res.Code.HasOpcode(isa.TailCall) {
		res.flags |= NoInline
	}
	return res
}

func (p *Procedure) HasFlag(f ProcFlag) bool {
	return p.flags&f != 0
}

// Call allows an abstract machine Procedure to be called for the purpose of
// instantiating a Closure. Only the compiler invokes this calling interface.
func (p *Procedure) Call(values ...data.Value) data.Value {
	return &Closure{
		Procedure: p,
		captured:  values,
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

// Equal compares this Procedure to another for equality
func (p *Procedure) Equal(other data.Value) bool {
	if p == other {
		return true
	}
	if other, ok := other.(*Procedure); ok {
		return p.Globals == other.Globals &&
			p.StackSize == other.StackSize &&
			p.LocalCount == other.LocalCount &&
			slices.Equal(p.Code, other.Code) &&
			p.Constants.Equal(other.Constants)
	}
	return false
}

func (p *Procedure) HashCode() uint64 {
	if h := atomic.LoadUint64(&p.hash); h != 0 {
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
	atomic.StoreUint64(&p.hash, res)
	return res
}

func (p *Procedure) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(p).Get(key)
}

// Captured returns the captured values of a Closure
func (c *Closure) Captured() data.Vector {
	return c.captured
}

// Call turns Closure into a Procedure, and serves as the virtual machine
func (c *Closure) Call(args ...data.Value) data.Value {
	return (&VM{CL: c, ARGS: args}).Run()
}

// CheckArity performs a compile-time arity check for the Closure
func (c *Closure) CheckArity(i int) error {
	return c.ArityChecker(i)
}

func (c *Closure) Equal(other data.Value) bool {
	if c == other {
		return true
	}
	if other, ok := other.(*Closure); ok {
		return c.Procedure.Equal(other.Procedure) &&
			c.captured.Equal(other.captured)
	}
	return false
}

func (c *Closure) HashCode() uint64 {
	if h := atomic.LoadUint64(&c.hash); h != 0 {
		return h
	}
	res := c.Procedure.HashCode()
	for i, v := range c.captured {
		res ^= data.HashCode(v)
		res ^= data.HashInt(i)
	}
	atomic.StoreUint64(&c.hash, res)
	return res
}
