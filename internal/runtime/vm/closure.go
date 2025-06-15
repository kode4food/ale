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
		hash     atomic.Uint64
	}

	argStack struct {
		prev *argStack
		args data.Vector
	}
)

// Error messages
const (
	// ErrBadInstruction is raised when the VM encounters an Opcode that has
	// not been properly mapped
	ErrBadInstruction = "unknown instruction encountered: %s"

	// ErrEmptyArgStack is raised when the VM encounters an instruction to pop
	// the argument stack, but the head of the stack is empty
	ErrEmptyArgStack = "attempt to pop empty argument stack"

	// ErrUnexpectedLabel is raised when the VM encounters a label that should
	// have otherwise been stripped when the instructions were made Runnable
	ErrUnexpectedLabel = "unexpected label encountered: %d"

	// ErrUnexpectedNoOp is raised when the VM encounters a NoOp instruction
	// that should have been stripped when the instructions were made Runnable
	ErrUnexpectedNoOp = "unexpected no-op encountered"
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

CurrentPC:
	INST = CODE[PC]
	switch INST.Opcode() {

	case isa.Add:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Add(MEM[SP].(data.Number))

	case isa.Arg:
		MEM[SP] = args[INST.Operand()]
		SP--

	case isa.ArgLen:
		MEM[SP] = data.Integer(len(args))
		SP--

	case isa.Bind:
		SP1 := SP + 1
		SP += 2
		name := MEM[SP1].(data.Local)
		value := MEM[SP]
		if err := bindOrShadow(c.Globals, name, value); err != nil {
			panic(err)
		}

	case isa.BindRef:
		SP1 := SP + 1
		SP += 2
		ref := MEM[SP1].(*Ref)
		ref.Value = MEM[SP]

	case isa.Call:
		op := INST.Operand()
		SP1 := SP + 1
		SP2 := SP1 + 1
		fn := MEM[SP1].(data.Procedure)
		callArgs := MEM[SP2 : SP2+int(op)]
		RES := SP1 + int(op)
		MEM[RES] = fn.Call(callArgs...)
		SP = RES - 1

	case isa.Call0:
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Procedure).Call()

	case isa.Call1:
		SP2 := SP + 2
		SP++
		MEM[SP2] = MEM[SP].(data.Procedure).Call(MEM[SP2])

	case isa.Call2:
		SP1 := SP + 1
		SP3 := SP + 3
		SP += 2
		MEM[SP3] = MEM[SP1].(data.Procedure).Call(MEM[SP], MEM[SP3])

	case isa.Call3:
		SP1 := SP + 1
		SP2 := SP + 2
		SP4 := SP + 4
		SP += 3
		MEM[SP4] = MEM[SP1].(data.Procedure).Call(MEM[SP2], MEM[SP], MEM[SP4])

	case isa.CallSelf:
		op := INST.Operand()
		SP1 := SP + 1
		callArgs := MEM[SP1 : SP1+int(op)]
		RES := SP + int(op)
		MEM[RES] = c.Call(callArgs...)
		SP = RES - 1

	case isa.CallWith:
		SP1 := SP + 2
		SP++
		callArgs := sequence.ToVector(MEM[SP1].(data.Sequence))
		MEM[SP1] = MEM[SP].(data.Procedure).Call(callArgs...)

	case isa.Car:
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Pair).Car()

	case isa.Cdr:
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Pair).Cdr()

	case isa.Closure:
		MEM[SP] = c.captured[INST.Operand()]
		SP--

	case isa.CondJump:
		SP++
		if MEM[SP] != data.False {
			PC = int(INST.Operand())
			goto CurrentPC
		}

	case isa.Cons:
		SP++
		SP1 := SP + 1
		if p, ok := MEM[SP1].(data.Prepender); ok {
			MEM[SP1] = p.Prepend(MEM[SP])
		} else {
			MEM[SP1] = data.NewCons(MEM[SP], MEM[SP1])
		}

	case isa.Const:
		MEM[SP] = c.Constants[INST.Operand()]
		SP--

	case isa.Deref:
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(*Ref).Value

	case isa.Div:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Div(MEM[SP].(data.Number))

	case isa.Dup:
		MEM[SP] = MEM[SP+1]
		SP--

	case isa.Empty:
		SP1 := SP + 1
		MEM[SP1] = data.Bool(MEM[SP1].(data.Sequence).IsEmpty())

	case isa.Eq:
		SP++
		SP1 := SP + 1
		MEM[SP1] = data.Bool(MEM[SP1].Equal(MEM[SP]))

	case isa.False:
		MEM[SP] = data.False
		SP--

	case isa.Jump:
		PC = int(INST.Operand())
		goto CurrentPC

	case isa.Label:
		// Labels should be stripped out by the compiler when made Runnable
		panic(debug.ProgrammerErrorf(ErrUnexpectedLabel, INST.Operand()))

	case isa.Load:
		MEM[SP] = MEM[LP+int(INST.Operand())]
		SP--

	case isa.Mod:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Mod(MEM[SP].(data.Number))

	case isa.Mul:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Mul(MEM[SP].(data.Number))

	case isa.Neg:
		SP1 := SP + 1
		MEM[SP1] = data.Integer(0).Sub(MEM[SP1].(data.Number))

	case isa.NegInt:
		MEM[SP] = -data.Integer(INST.Operand())
		SP--

	case isa.NewRef:
		MEM[SP] = new(Ref)
		SP--

	case isa.NoOp:
		// NoOp should be stripped out by the compiler when made Runnable
		panic(debug.ProgrammerError(ErrUnexpectedNoOp))

	case isa.Not:
		SP1 := SP + 1
		MEM[SP1] = !MEM[SP1].(data.Bool)

	case isa.Null:
		MEM[SP] = data.Null
		SP--

	case isa.NumEq:
		SP++
		SP1 := SP + 1
		cmp := MEM[SP1].(data.Number).Cmp(MEM[SP].(data.Number))
		MEM[SP1] = data.Bool(data.EqualTo == cmp)

	case isa.NumGt:
		SP++
		SP1 := SP + 1
		cmp := MEM[SP1].(data.Number).Cmp(MEM[SP].(data.Number))
		MEM[SP1] = data.Bool(data.GreaterThan == cmp)

	case isa.NumGte:
		SP++
		SP1 := SP + 1
		cmp := MEM[SP1].(data.Number).Cmp(MEM[SP].(data.Number))
		MEM[SP1] = data.Bool(cmp == data.GreaterThan || cmp == data.EqualTo)

	case isa.NumLt:
		SP++
		SP1 := SP + 1
		cmp := MEM[SP1].(data.Number).Cmp(MEM[SP].(data.Number))
		MEM[SP1] = data.Bool(data.LessThan == cmp)

	case isa.NumLte:
		SP++
		SP1 := SP + 1
		cmp := MEM[SP1].(data.Number).Cmp(MEM[SP].(data.Number))
		MEM[SP1] = data.Bool(cmp == data.LessThan || cmp == data.EqualTo)

	case isa.Panic:
		panic(errors.New(data.ToString(MEM[SP+1])))

	case isa.Pop:
		SP++

	case isa.PopArgs:
		if AP == nil {
			panic(debug.ProgrammerError(ErrEmptyArgStack))
		}
		args = AP.args
		AP = AP.prev

	case isa.PosInt:
		MEM[SP] = data.Integer(INST.Operand())
		SP--

	case isa.Private:
		SP++
		if _, err := c.Globals.Private(MEM[SP].(data.Local)); err != nil {
			panic(err)
		}

	case isa.Public:
		SP++
		if _, err := c.Globals.Public(MEM[SP].(data.Local)); err != nil {
			panic(err)
		}

	case isa.PushArgs:
		RES := SP + int(INST.Operand())
		AP = &argStack{
			args: args,
			prev: AP,
		}
		args = slices.Clone(MEM[SP+1 : RES+1])
		SP = RES

	case isa.Resolve:
		SP1 := SP + 1
		MEM[SP1] = env.MustResolveValue(c.Globals, MEM[SP1].(data.Symbol))

	case isa.RestArg:
		MEM[SP] = data.Vector(args[INST.Operand():])
		SP--

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

	case isa.Sub:
		SP++
		SP1 := SP + 1
		MEM[SP1] = MEM[SP1].(data.Number).Sub(MEM[SP].(data.Number))

	case isa.TailCall: // Fully dynamic tail call
		op := INST.Operand()
		SP1 := SP + 1
		SP2 := SP1 + 1
		val := MEM[SP1]
		callArgs := MEM[SP2 : SP2+int(op)]
		cl, ok := val.(*Closure)
		if !ok {
			return val.(data.Procedure).Call(callArgs...)
		}
		args = slices.Clone(callArgs)
		if cl == c {
			goto InitState
		}
		c = cl
		if len(MEM) < int(c.StackSize+c.LocalCount) {
			free(MEM)
			goto InitMem
		}
		goto InitCode

	case isa.TailClos:
		op := INST.Operand()
		SP1 := SP + 1
		SP2 := SP1 + 1
		c = MEM[SP1].(*Closure)
		args = slices.Clone(MEM[SP2 : SP2+int(op)])
		if len(MEM) < int(c.StackSize+c.LocalCount) {
			free(MEM)
			goto InitMem
		}
		goto InitCode

	case isa.TailSelf:
		op := INST.Operand()
		SP1 := SP + 1
		args = slices.Clone(MEM[SP1 : SP1+int(op)])
		goto InitState

	case isa.True:
		MEM[SP] = data.True
		SP--

	case isa.Vector:
		op := INST.Operand()
		RES := SP + int(op)
		MEM[RES] = slices.Clone(MEM[SP+1 : RES+1])
		SP = RES - 1

	case isa.Zero:
		MEM[SP] = data.Integer(0)
		SP--

	default:
		panic(debug.ProgrammerErrorf(ErrBadInstruction, INST))
	}

	PC++
	goto CurrentPC
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
	if h := c.hash.Load(); h != 0 {
		return h
	}
	res := c.Procedure.HashCode()
	for i, v := range c.captured {
		res ^= data.HashCode(v)
		res ^= data.HashInt(i)
	}
	c.hash.Store(res)
	return res
}

func bindOrShadow(ns env.Namespace, n data.Local, v data.Value) error {
	e, in, err := ns.Resolve(n)
	if err != nil || in != ns {
		return env.BindPublic(ns, n, v)
	}
	return e.Bind(v)
}
