package vm

import (
	"errors"
	"slices"
	"sync/atomic"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

type (
	Closure struct {
		*Procedure
		captured data.Vector
		hash     uint64
	}

	argStack struct {
		prev *argStack
		args data.Vector
	}
)

// Error messages
const (
	// ErrBadInstruction is raised when the VM encounters an Opcode that has not
	// been properly mapped
	ErrBadInstruction = "unknown instruction encountered: %s"

	// ErrEmptyArgStack is raised when the VM encounters an instruction to pop
	// the argument stack, but the head of the stack is empty
	ErrEmptyArgStack = "attempt to pop empty argument stack"
)

// Captured returns the captured values of a Closure
func (c *Closure) Captured() data.Vector {
	return c.captured
}

// Call turns Closure into a Procedure, and serves as the virtual machine
func (c *Closure) Call(args ...data.Value) data.Value {
	var MEM data.Vector
	var CODE isa.Instructions
	var PC, LP, SP int
	var INST isa.Instruction
	var AP *argStack

	defer func() { free(MEM) }()

InitMem:
	MEM = malloc(int(c.StackSize + c.LocalCount))

InitCode:
	CODE = c.Code
	LP = int(c.StackSize)

InitState:
	SP = LP - 1
	PC = 0
	goto CurrentPC

NextPC:
	PC++

CurrentPC:
	INST = CODE[PC]
	switch INST.Opcode() {

	case isa.Add:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Add(
			MEM[SP].(data.Number),
		)
		goto NextPC

	case isa.Arg:
		MEM[SP] = args[INST.Operand()]
		SP--
		goto NextPC

	case isa.ArgLen:
		MEM[SP] = data.Integer(len(args))
		SP--
		goto NextPC

	case isa.Bind:
		SP++
		name := MEM[SP].(data.Local)
		SP++
		if err := c.Globals.Declare(name).Bind(MEM[SP]); err != nil {
			panic(err)
		}
		goto NextPC

	case isa.BindRef:
		SP++
		ref := MEM[SP].(*Ref)
		SP++
		ref.Value = MEM[SP]
		goto NextPC

	case isa.Call0:
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Procedure).Call()
		goto NextPC

	case isa.Call1:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP].(data.Procedure).Call(MEM[SP1])
		goto NextPC

	case isa.Call:
		op := INST.Operand()
		SP1 := SP + 1
		SP2 := SP1 + 1
		fn := MEM[SP1].(data.Procedure)
		args := slices.Clone(MEM[SP2 : SP2+int(op)])
		RES := SP1 + int(op)
		MEM[RES] = fn.Call(args...)
		SP = RES - 1
		goto NextPC

	case isa.CallWith:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP].(data.Procedure).Call(
			sequence.ToValues(MEM[SP1].(data.Sequence))...,
		)
		goto NextPC

	case isa.Car:
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Pair).Car()
		goto NextPC

	case isa.Cdr:
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Pair).Cdr()
		goto NextPC

	case isa.Closure:
		MEM[SP] = c.captured[INST.Operand()]
		SP--
		goto NextPC

	case isa.CondJump:
		SP++
		if MEM[SP] != data.False {
			PC = int(INST.Operand())
			goto CurrentPC
		}
		goto NextPC

	case isa.Cons:
		SP++
		SP1 := SP + 1
		if p, ok := MEM[SP1].(data.Prepender); ok {
			MEM[SP1] = p.Prepend(MEM[SP])
			goto NextPC
		}
		MEM[SP1] = data.NewCons(MEM[SP], MEM[SP1])
		goto NextPC

	case isa.Const:
		MEM[SP] = c.Constants[INST.Operand()]
		SP--
		goto NextPC

	case isa.Declare:
		SP++
		c.Globals.Declare(
			MEM[SP].(data.Local),
		)
		goto NextPC

	case isa.Deref:
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(*Ref).Value
		goto NextPC

	case isa.Div:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Div(
			MEM[SP].(data.Number),
		)
		goto NextPC

	case isa.Dup:
		MEM[SP] = MEM[SP+1]
		SP--
		goto NextPC

	case isa.Empty:
		SP1 := SP + 1
		MEM[SP1] = data.Bool(MEM[SP1].(data.Sequence).IsEmpty())
		goto NextPC

	case isa.Eq:
		SP++
		SP1 := SP + 1
		MEM[SP1] = data.Bool(MEM[SP1].Equal(MEM[SP]))
		goto NextPC

	case isa.False:
		MEM[SP] = data.False
		SP--
		goto NextPC

	case isa.Jump:
		PC = int(INST.Operand())
		goto CurrentPC

	case isa.Load:
		MEM[SP] = MEM[LP+int(INST.Operand())]
		SP--
		goto NextPC

	case isa.Mod:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Mod(
			MEM[SP].(data.Number),
		)
		goto NextPC

	case isa.Mul:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Mul(
			MEM[SP].(data.Number),
		)
		goto NextPC

	case isa.Neg:
		SP1 := SP + 1
		MEM[SP1] = data.Integer(0).Sub(
			MEM[SP1].(data.Number),
		)
		goto NextPC

	case isa.NegInt:
		MEM[SP] = -data.Integer(INST.Operand())
		SP--
		goto NextPC

	case isa.NewRef:
		MEM[SP] = new(Ref)
		SP--
		goto NextPC

	case isa.NoOp:
		goto NextPC

	case isa.Not:
		SP1 := SP + 1
		MEM[SP1] = !MEM[SP1].(data.Bool)
		goto NextPC

	case isa.Null:
		MEM[SP] = data.Null
		SP--
		goto NextPC

	case isa.NumEq:
		SP++
		SP1 := SP + 1
		MEM[SP1] = data.Bool(
			data.EqualTo == MEM[SP1].(data.Number).Cmp(
				MEM[SP].(data.Number),
			),
		)
		goto NextPC

	case isa.NumGt:
		SP++
		SP1 := SP + 1
		MEM[SP1] = data.Bool(
			data.GreaterThan == MEM[SP1].(data.Number).Cmp(
				MEM[SP].(data.Number),
			),
		)
		goto NextPC

	case isa.NumGte:
		SP++
		SP1 := SP + 1
		cmp := MEM[SP1].(data.Number).Cmp(
			MEM[SP].(data.Number),
		)
		MEM[SP1] = data.Bool(
			cmp == data.GreaterThan || cmp == data.EqualTo,
		)
		goto NextPC

	case isa.NumLt:
		SP++
		SP1 := SP + 1
		MEM[SP1] = data.Bool(
			data.LessThan == MEM[SP1].(data.Number).Cmp(
				MEM[SP].(data.Number),
			),
		)
		goto NextPC

	case isa.NumLte:
		SP++
		SP1 := SP + 1
		cmp := MEM[SP1].(data.Number).Cmp(
			MEM[SP].(data.Number),
		)
		MEM[SP1] = data.Bool(
			cmp == data.LessThan || cmp == data.EqualTo,
		)
		goto NextPC

	case isa.Panic:
		panic(errors.New(data.ToString(MEM[SP+1])))

	case isa.Pop:
		SP++
		goto NextPC

	case isa.PopArgs:
		if AP == nil {
			panic(debug.ProgrammerError(ErrEmptyArgStack))
		}
		args = AP.args
		AP = AP.prev
		goto NextPC

	case isa.PosInt:
		MEM[SP] = data.Integer(INST.Operand())
		SP--
		goto NextPC

	case isa.Private:
		SP++
		c.Globals.Private(
			MEM[SP].(data.Local),
		)
		goto NextPC

	case isa.PushArgs:
		RES := SP + int(INST.Operand())
		AP = &argStack{
			args: args,
			prev: AP,
		}
		args = slices.Clone(MEM[SP+1 : RES+1])
		SP = RES
		goto NextPC

	case isa.Resolve:
		SP1 := SP + 1
		MEM[SP1] = env.MustResolveValue(
			c.Globals,
			MEM[SP1].(data.Symbol),
		)
		goto NextPC

	case isa.RestArg:
		MEM[SP] = data.Vector(args[INST.Operand():])
		SP--
		goto NextPC

	case isa.RetFalse:
		return data.False

	case isa.RetNull:
		return data.Null

	case isa.RetTrue:
		return data.True

	case isa.Return:
		return MEM[SP+1]

	case isa.Store:
		SP++
		MEM[LP+int(INST.Operand())] = MEM[SP]
		goto NextPC

	case isa.Sub:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Sub(
			MEM[SP].(data.Number),
		)
		goto NextPC

	case isa.TailCall:
		op := INST.Operand()
		SP1 := SP + 1
		SP2 := SP1 + 1
		val := MEM[SP1]
		args = slices.Clone(MEM[SP2 : SP2+int(op)])
		cl, ok := val.(*Closure)
		if !ok {
			return val.(data.Procedure).Call(args...)
		}
		if cl == c {
			goto InitState
		}
		c = cl // intentional
		if len(MEM) < int(c.StackSize+c.LocalCount) {
			free(MEM)
			goto InitMem
		}
		goto InitCode

	case isa.True:
		MEM[SP] = data.True
		SP--
		goto NextPC

	case isa.Vector:
		op := INST.Operand()
		RES := SP + int(op)
		MEM[RES] = slices.Clone(MEM[SP+1 : RES+1])
		SP = RES - 1
		goto NextPC

	case isa.Zero:
		MEM[SP] = data.Integer(0)
		SP--
		goto NextPC

	default:
		panic(debug.ProgrammerError(ErrBadInstruction, INST))

	}
}

// CheckArity performs a compile-time arity check for the Closure
func (c *Closure) CheckArity(i int) error {
	return c.ArityChecker(i)
}

func (c *Closure) Equal(other data.Value) bool {
	if other, ok := other.(*Closure); ok {
		return c == other ||
			c.Procedure.Equal(other.Procedure) &&
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
