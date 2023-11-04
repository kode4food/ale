package vm

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/runtime/isa"
)

type closure struct {
	*Procedure
	values data.Values
}

func newClosure(lambda *Procedure, values data.Values) *closure {
	return &closure{
		Procedure: lambda,
		values:    values,
	}
}

// Call turns closure into a Procedure, and serves as the virtual machine
func (c *closure) Call(args ...data.Value) data.Value {
	var (
		CODE isa.Instructions
		MEM  data.Values
		PC   unsafe.Pointer
		LP   int
		SP   int
	)

initMem:
	MEM = make(data.Values, c.Procedure.StackSize+c.Procedure.LocalCount)

initCode:
	CODE = c.Procedure.Code
	LP = c.Procedure.StackSize

initState:
	SP = LP - 1
	// cheaper than a goto
	PC = unsafe.Add(unsafe.Pointer(&CODE[0]), -int(unsafe.Sizeof(CODE[0])))

nextPC:
	PC = unsafe.Add(PC, unsafe.Sizeof(CODE[0]))

opSwitch:
	oc, op := (*(*isa.Instruction)(PC)).Split()
	switch oc {
	case isa.Null:
		MEM[SP] = data.Null
		SP--
		goto nextPC

	case isa.Zero:
		MEM[SP] = data.Integer(0)
		SP--
		goto nextPC

	case isa.PosInt:
		MEM[SP] = data.Integer(op)
		SP--
		goto nextPC

	case isa.NegInt:
		MEM[SP] = -data.Integer(op)
		SP--
		goto nextPC

	case isa.True:
		MEM[SP] = data.True
		SP--
		goto nextPC

	case isa.False:
		MEM[SP] = data.False
		SP--
		goto nextPC

	case isa.Const:
		MEM[SP] = c.Procedure.Constants[op]
		SP--
		goto nextPC

	case isa.Arg:
		MEM[SP] = args[op]
		SP--
		goto nextPC

	case isa.RestArg:
		MEM[SP] = data.NewVector(args[op:]...)
		SP--
		goto nextPC

	case isa.ArgLen:
		MEM[SP] = data.Integer(len(args))
		SP--
		goto nextPC

	case isa.Closure:
		MEM[SP] = c.values[op]
		SP--
		goto nextPC

	case isa.Load:
		MEM[SP] = MEM[LP+int(op)]
		SP--
		goto nextPC

	case isa.Store:
		SP++
		MEM[LP+int(op)] = MEM[SP]
		goto nextPC

	case isa.NewRef:
		MEM[SP] = new(Ref)
		SP--
		goto nextPC

	case isa.BindRef:
		SP++
		ref := MEM[SP].(*Ref)
		SP++
		ref.Value = MEM[SP]
		goto nextPC

	case isa.Deref:
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(*Ref).Value
		goto nextPC

	case isa.Car:
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(data.Pair).Car()
		goto nextPC

	case isa.Cdr:
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(data.Pair).Cdr()
		goto nextPC

	case isa.Cons:
		SP++
		SP1 := &MEM[SP+1]
		if p, ok := (*SP1).(data.Prepender); ok {
			*SP1 = p.Prepend(MEM[SP])
			goto nextPC
		}
		*SP1 = data.NewCons(MEM[SP], *SP1)
		goto nextPC

	case isa.Empty:
		SP1 := &MEM[SP+1]
		*SP1 = data.Bool((*SP1).(data.Sequence).IsEmpty())
		goto nextPC

	case isa.Eq:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = data.Bool((*SP1).Equal(MEM[SP]))
		goto nextPC

	case isa.Not:
		SP1 := &MEM[SP+1]
		*SP1 = !(*SP1).(data.Bool)
		goto nextPC

	case isa.Declare:
		SP++
		c.Procedure.Globals.Declare(
			MEM[SP].(data.Local),
		)
		goto nextPC

	case isa.Private:
		SP++
		c.Procedure.Globals.Private(
			MEM[SP].(data.Local),
		)
		goto nextPC

	case isa.Bind:
		SP++
		name := MEM[SP].(data.Local)
		SP++
		c.Procedure.Globals.Declare(name).Bind(MEM[SP])
		goto nextPC

	case isa.Resolve:
		SP1 := &MEM[SP+1]
		*SP1 = env.MustResolveValue(
			c.Procedure.Globals,
			(*SP1).(data.Symbol),
		)
		goto nextPC

	case isa.Dup:
		MEM[SP] = MEM[SP+1]
		SP--
		goto nextPC

	case isa.Pop:
		SP++
		goto nextPC

	case isa.Add:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(data.Number).Add(
			MEM[SP].(data.Number),
		)
		goto nextPC

	case isa.Sub:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(data.Number).Sub(
			MEM[SP].(data.Number),
		)
		goto nextPC

	case isa.Mul:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(data.Number).Mul(
			MEM[SP].(data.Number),
		)
		goto nextPC

	case isa.Div:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(data.Number).Div(
			MEM[SP].(data.Number),
		)
		goto nextPC

	case isa.Mod:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(data.Number).Mod(
			MEM[SP].(data.Number),
		)
		goto nextPC

	case isa.NumEq:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = data.Bool(
			data.EqualTo == (*SP1).(data.Number).Cmp(
				MEM[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.NumLt:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = data.Bool(
			data.LessThan == (*SP1).(data.Number).Cmp(
				MEM[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.NumLte:
		SP++
		SP1 := &MEM[SP+1]
		cmp := (*SP1).(data.Number).Cmp(
			MEM[SP].(data.Number),
		)
		*SP1 = data.Bool(
			cmp == data.LessThan || cmp == data.EqualTo,
		)
		goto nextPC

	case isa.NumGt:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = data.Bool(
			data.GreaterThan == (*SP1).(data.Number).Cmp(
				MEM[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.NumGte:
		SP++
		SP1 := &MEM[SP+1]
		cmp := (*SP1).(data.Number).Cmp(
			MEM[SP].(data.Number),
		)
		*SP1 = data.Bool(
			cmp == data.GreaterThan || cmp == data.EqualTo,
		)
		goto nextPC

	case isa.Neg:
		SP1 := &MEM[SP+1]
		*SP1 = data.Integer(0).Sub(
			(*SP1).(data.Number),
		)
		goto nextPC

	case isa.Call0:
		SP1 := &MEM[SP+1]
		*SP1 = (*SP1).(data.Procedure).Call()
		goto nextPC

	case isa.Call1:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = MEM[SP].(data.Procedure).Call(*SP1)
		goto nextPC

	case isa.Call:
		SP1 := SP + 1
		fn := MEM[SP1].(data.Procedure)
		// prepare args
		args := make(data.Values, op)
		copy(args, MEM[SP1+1:LP]) // because stack mutates
		// call function
		RES := SP1 + int(op)
		MEM[RES] = fn.Call(args...)
		SP = RES - 1
		goto nextPC

	case isa.CallWith:
		SP++
		SP1 := &MEM[SP+1]
		*SP1 = MEM[SP].(data.Procedure).Call(
			sequence.ToValues((*SP1).(data.Sequence))...,
		)
		goto nextPC

	case isa.TailCall:
		SP1 := SP + 1
		val := MEM[SP1]
		// prepare args
		args = make(data.Values, op)
		copy(args, MEM[SP1+1:LP]) // because stack mutates
		// call function
		cl, ok := val.(*closure)
		if !ok {
			return val.(data.Procedure).Call(args...)
		}
		if cl == c {
			goto initState
		}
		c = cl // intentional
		ss := c.Procedure.StackSize
		lc := c.Procedure.LocalCount
		if len(MEM) < ss+lc {
			goto initMem
		}
		goto initCode

	case isa.Jump:
		PC = unsafe.Pointer(&CODE[int(op)])
		goto opSwitch

	case isa.CondJump:
		SP++
		if MEM[SP] != data.False {
			PC = unsafe.Pointer(&CODE[int(op)])
			goto opSwitch
		}
		goto nextPC

	case isa.NoOp:
		goto nextPC

	case isa.Panic:
		panic(errors.New(MEM[SP+1].String()))

	case isa.Return:
		return MEM[SP+1]

	case isa.RetNull:
		return data.Null

	case isa.RetTrue:
		return data.True

	case isa.RetFalse:
		return data.False

	default:
		// Programmer error
		panic(fmt.Sprintf("unknown opcode: %s", oc))
	}
}

// CheckArity performs a compile-time arity check for the closure
func (c *closure) CheckArity(i int) error {
	return c.ArityChecker(i)
}

func (c *closure) Equal(v data.Value) bool {
	return c == v
}
