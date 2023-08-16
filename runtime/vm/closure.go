package vm

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	errUnknownOpcode = "unknown opcode: %s"
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
		CODE   []isa.Word
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
	switch isa.Opcode(CODE[PC]) {
	case isa.Nil:
		STACK[SP] = data.Nil
		SP--
		goto nextPC

	case isa.Zero:
		STACK[SP] = data.Integer(0)
		SP--
		goto nextPC

	case isa.One:
		STACK[SP] = data.Integer(1)
		SP--
		goto nextPC

	case isa.NegOne:
		STACK[SP] = data.Integer(-1)
		SP--
		goto nextPC

	case isa.Two:
		STACK[SP] = data.Integer(2)
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
		PC++
		STACK[SP] = c.Lambda.Constants[CODE[PC]]
		SP--
		goto nextPC

	case isa.Arg:
		PC++
		STACK[SP] = args[CODE[PC]]
		SP--
		goto nextPC

	case isa.RestArg:
		PC++
		STACK[SP] = data.NewVector(args[CODE[PC]:]...)
		SP--
		goto nextPC

	case isa.ArgLen:
		STACK[SP] = data.Integer(len(args))
		SP--
		goto nextPC

	case isa.Closure:
		PC++
		STACK[SP] = c.values[CODE[PC]]
		SP--
		goto nextPC

	case isa.Load:
		PC++
		STACK[SP] = LOCALS[CODE[PC]]
		SP--
		goto nextPC

	case isa.Store:
		PC++
		SP++
		LOCALS[CODE[PC]] = STACK[SP]
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
		SP1 := SP + 1
		STACK[SP1] = STACK[SP1].(*Ref).Value
		goto nextPC

	case isa.Declare:
		SP++
		c.Lambda.Globals.Declare(
			STACK[SP].(data.Name),
		)
		goto nextPC

	case isa.Private:
		SP++
		c.Lambda.Globals.Private(
			STACK[SP].(data.Name),
		)
		goto nextPC

	case isa.Bind:
		SP++
		name := STACK[SP].(data.Name)
		SP++
		c.Lambda.Globals.Declare(name).Bind(STACK[SP])
		goto nextPC

	case isa.Resolve:
		SP1 := SP + 1
		STACK[SP1] = env.MustResolveValue(
			c.Lambda.Globals,
			STACK[SP1].(data.Symbol),
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
		SP1 := SP + 1
		STACK[SP1] = STACK[SP1].(data.Number).Add(
			STACK[SP].(data.Number),
		)
		goto nextPC

	case isa.Sub:
		SP++
		SP1 := SP + 1
		STACK[SP1] = STACK[SP1].(data.Number).Sub(
			STACK[SP].(data.Number),
		)
		goto nextPC

	case isa.Mul:
		SP++
		SP1 := SP + 1
		STACK[SP1] = STACK[SP1].(data.Number).Mul(
			STACK[SP].(data.Number),
		)
		goto nextPC

	case isa.Div:
		SP++
		SP1 := SP + 1
		STACK[SP1] = STACK[SP1].(data.Number).Div(
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
		SP1 := SP + 1
		STACK[SP1] = data.Bool(
			data.EqualTo == STACK[SP1].(data.Number).Cmp(
				STACK[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.Neq:
		SP++
		SP1 := SP + 1
		STACK[SP1] = data.Bool(
			data.EqualTo != STACK[SP1].(data.Number).Cmp(
				STACK[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.Lt:
		SP++
		SP1 := SP + 1
		STACK[SP1] = data.Bool(
			data.LessThan == STACK[SP1].(data.Number).Cmp(
				STACK[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.Lte:
		SP++
		SP1 := SP + 1
		cmp := STACK[SP1].(data.Number).Cmp(
			STACK[SP].(data.Number),
		)
		STACK[SP1] = data.Bool(
			cmp == data.LessThan || cmp == data.EqualTo,
		)
		goto nextPC

	case isa.Gt:
		SP++
		SP1 := SP + 1
		STACK[SP1] = data.Bool(
			data.GreaterThan == STACK[SP1].(data.Number).Cmp(
				STACK[SP].(data.Number),
			),
		)
		goto nextPC

	case isa.Gte:
		SP++
		SP1 := SP + 1
		cmp := STACK[SP1].(data.Number).Cmp(
			STACK[SP].(data.Number),
		)
		STACK[SP1] = data.Bool(
			cmp == data.GreaterThan || cmp == data.EqualTo,
		)
		goto nextPC

	case isa.Neg:
		SP1 := SP + 1
		STACK[SP1] = data.Integer(0).Sub(
			STACK[SP1].(data.Number),
		)
		goto nextPC

	case isa.Not:
		SP1 := SP + 1
		STACK[SP1] = !STACK[SP1].(data.Bool)
		goto nextPC

	case isa.MakeTruthy:
		SP1 := SP + 1
		STACK[SP1] = data.Bool(
			data.Truthy(STACK[SP1]),
		)
		goto nextPC

	case isa.Call0:
		SP1 := SP + 1
		STACK[SP1] = STACK[SP1].(data.Function).Call()
		goto nextPC

	case isa.Call1:
		SP++
		SP1 := SP + 1
		STACK[SP1] = STACK[SP].(data.Function).Call(STACK[SP1])
		goto nextPC

	case isa.Call:
		PC++
		SP1 := SP + 1
		// prepare args
		argCount := int(CODE[PC])
		args := make(data.Values, argCount)
		copy(args, STACK[SP1+1:]) // must be a copy
		// call function
		RES := SP1 + argCount
		STACK[RES] = STACK[SP1].(data.Function).Call(args...)
		SP = RES - 1
		goto nextPC

	case isa.TailCall:
		SP1 := SP + 1
		// prepare args
		argCount := int(CODE[PC+1])
		args = make(data.Values, argCount)
		copy(args, STACK[SP1+1:]) // must be a copy
		// call function
		val := STACK[SP1]
		if vc, ok := val.(*closure); ok {
			if vc == c {
				goto initState
			}
			c = vc // intentional
			if len(DATA) < c.Lambda.StackSize+c.Lambda.LocalCount {
				goto initFrame
			}
			CODE = c.Lambda.Code
			if len(STACK) != c.Lambda.StackSize {
				STACK = DATA[0:c.Lambda.StackSize]
			}
			if len(LOCALS) != c.Lambda.LocalCount {
				LOCALS = DATA[len(DATA)-c.Lambda.LocalCount:]
			}
			goto initState
		}
		return val.(data.Function).Call(args...)

	case isa.Jump:
		PC = int(CODE[PC+1])
		goto opSwitch

	case isa.CondJump:
		SP++
		if STACK[SP].(data.Bool) {
			PC = int(CODE[PC+1])
			goto opSwitch
		}
		PC += 2
		goto opSwitch

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
		panic(fmt.Errorf(errUnknownOpcode, isa.Opcode(CODE[PC])))
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
