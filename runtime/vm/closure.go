package vm

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/runtime/isa"
)

type closure struct {
	*Lambda
	values data.Values
}

func newClosure(lambda *Lambda, values data.Values) *closure {
	return &closure{
		Lambda: lambda,
		values: values,
	}
}

// Call turns closure into a Function, and serves as the virtual machine
func (c *closure) Call(args ...data.Value) data.Value {
	var (
		DATA   data.Values
		CODE   isa.Instructions
		STACK  data.Values
		LOCALS data.Values
		SP     int
		PC     int
	)

initFrame:
	CODE = c.Lambda.Code
	DATA = make(data.Values, c.Lambda.StackSize+c.Lambda.LocalCount)
	STACK = DATA[0:c.Lambda.StackSize]
	LOCALS = DATA[c.Lambda.StackSize:]

initState:
	SP = len(STACK) - 1
	PC = -1 // cheaper than a goto

nextPC:
	PC++

opSwitch:
	oc, op := CODE[PC].Split()
	switch oc {
	case isa.Nil:
		STACK[SP] = data.Nil
		SP--
		goto nextPC

	case isa.Zero:
		STACK[SP] = data.Integer(0)
		SP--
		goto nextPC

	case isa.PosInt:
		STACK[SP] = data.Integer(op)
		SP--
		goto nextPC

	case isa.NegInt:
		STACK[SP] = -data.Integer(op)
		SP--
		goto nextPC

	case isa.True:
		STACK[SP] = data.True
		SP--
		goto nextPC

	case isa.False:
		STACK[SP] = data.False
		SP--
		goto nextPC

	case isa.Const:
		STACK[SP] = c.Lambda.Constants[op]
		SP--
		goto nextPC

	case isa.Arg:
		STACK[SP] = args[op]
		SP--
		goto nextPC

	case isa.RestArg:
		STACK[SP] = data.NewVector(args[op:]...)
		SP--
		goto nextPC

	case isa.ArgLen:
		STACK[SP] = data.Integer(len(args))
		SP--
		goto nextPC

	case isa.Closure:
		STACK[SP] = c.values[op]
		SP--
		goto nextPC

	case isa.Load:
		STACK[SP] = LOCALS[op]
		SP--
		goto nextPC

	case isa.Store:
		SP++
		LOCALS[op] = STACK[SP]
		goto nextPC

	case isa.NewRef:
		STACK[SP] = new(Ref)
		SP--
		goto nextPC

	case isa.BindRef:
		SP++
		ref := STACK[SP].(*Ref)
		SP++
		ref.Value = STACK[SP]
		goto nextPC

	case isa.Deref:
		SP1 := &STACK[SP+1]
		*SP1 = (*SP1).(*Ref).Value
		goto nextPC

	case isa.Car:
		SP1 := &STACK[SP+1]
		*SP1 = (*SP1).(data.Pair).Car()
		goto nextPC

	case isa.Cdr:
		SP1 := &STACK[SP+1]
		*SP1 = (*SP1).(data.Pair).Cdr()
		goto nextPC

	case isa.Declare:
		SP++
		c.Lambda.Globals.Declare(
			STACK[SP].(data.Local),
		)
		goto nextPC

	case isa.Private:
		SP++
		c.Lambda.Globals.Private(
			STACK[SP].(data.Local),
		)
		goto nextPC

	case isa.Bind:
		SP++
		val := STACK[SP]
		SP++
		c.Lambda.Globals.Declare(
			STACK[SP].(data.Local),
		).Bind(val)
		goto nextPC

	case isa.Resolve:
		SP1 := &STACK[SP+1]
		*SP1 = env.MustResolveValue(
			c.Lambda.Globals,
			(*SP1).(data.Symbol),
		)
		goto nextPC

	case isa.Dup:
		STACK[SP] = STACK[SP+1]
		SP--
		goto nextPC

	case isa.Pop:
		SP++
		goto nextPC

	case isa.Add:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = (*SP1).(data.Number).Add(
			STACK[SP].(data.Number),
		)
		goto nextPC

	case isa.Sub:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = (*SP1).(data.Number).Sub(
			STACK[SP].(data.Number),
		)
		goto nextPC

	case isa.Mul:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = (*SP1).(data.Number).Mul(
			STACK[SP].(data.Number),
		)
		goto nextPC

	case isa.Div:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = (*SP1).(data.Number).Div(
			STACK[SP].(data.Number),
		)
		goto nextPC

	case isa.Mod:
		SP++
		SP1 := SP + 1
		STACK[SP1] = STACK[SP1].(data.Number).Mod(
			STACK[SP].(data.Number),
		)
		goto nextPC

	case isa.Eq:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = data.Bool(
			data.EqualTo == (*SP1).(data.Number).Cmp(
				STACK[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.Neq:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = data.Bool(
			data.EqualTo != (*SP1).(data.Number).Cmp(
				STACK[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.Lt:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = data.Bool(
			data.LessThan == (*SP1).(data.Number).Cmp(
				STACK[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.Lte:
		SP++
		SP1 := &STACK[SP+1]
		cmp := (*SP1).(data.Number).Cmp(
			STACK[SP].(data.Number),
		)
		*SP1 = data.Bool(
			cmp == data.LessThan || cmp == data.EqualTo,
		)
		goto nextPC

	case isa.Gt:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = data.Bool(
			data.GreaterThan == (*SP1).(data.Number).Cmp(
				STACK[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.Gte:
		SP++
		SP1 := &STACK[SP+1]
		cmp := (*SP1).(data.Number).Cmp(
			STACK[SP].(data.Number),
		)
		*SP1 = data.Bool(
			cmp == data.GreaterThan || cmp == data.EqualTo,
		)
		goto nextPC

	case isa.Neg:
		SP1 := &STACK[SP+1]
		*SP1 = data.Integer(0).Sub(
			(*SP1).(data.Number),
		)
		goto nextPC

	case isa.Not:
		SP1 := &STACK[SP+1]
		*SP1 = !(*SP1).(data.Bool)
		goto nextPC

	case isa.MakeTruthy:
		SP1 := &STACK[SP+1]
		*SP1 = data.Bool(
			data.Truthy(*SP1),
		)
		goto nextPC

	case isa.Call0:
		SP1 := &STACK[SP+1]
		*SP1 = (*SP1).(data.Function).Call()
		goto nextPC

	case isa.Call1:
		SP++
		SP1 := &STACK[SP+1]
		*SP1 = STACK[SP].(data.Function).Call(*SP1)
		goto nextPC

	case isa.Call:
		SP1 := SP + 1
		fn := STACK[SP1].(data.Function)
		// prepare args
		args := make(data.Values, op)
		copy(args, STACK[SP1+1:]) // because stack mutates
		// call function
		RES := SP1 + int(op)
		STACK[RES] = fn.Call(args...)
		SP = RES - 1
		goto nextPC

	case isa.TailCall:
		SP1 := SP + 1
		val := STACK[SP1]
		// prepare args
		args = make(data.Values, op)
		copy(args, STACK[SP1+1:]) // because stack mutates
		// call function
		cl, ok := val.(*closure)
		if !ok {
			return val.(data.Function).Call(args...)
		}
		if cl == c {
			goto initState
		}
		c = cl // intentional
		ss := c.Lambda.StackSize
		lc := c.Lambda.LocalCount
		if len(DATA) < ss+lc {
			goto initFrame
		}
		CODE = c.Lambda.Code
		if len(STACK) != ss {
			STACK = DATA[0:ss]
		}
		if len(LOCALS) != lc {
			LOCALS = DATA[len(DATA)-lc:]
		}
		goto initState

	case isa.Jump:
		PC = int(op)
		goto opSwitch

	case isa.CondJump:
		SP++
		if STACK[SP].(data.Bool) {
			PC = int(op)
			goto opSwitch
		}
		goto nextPC

	case isa.Panic:
		panic(errors.New(STACK[SP+1].String()))

	case isa.Return:
		return STACK[SP+1]

	case isa.RetNil:
		return data.Nil

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

// Convention returns the closure's calling convention
func (c *closure) Convention() data.Convention {
	return data.ApplicativeCall
}

func (c *closure) Equal(v data.Value) bool {
	return c == v
}

func (c *closure) String() string {
	return data.DumpString(c)
}
