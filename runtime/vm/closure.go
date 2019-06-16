package vm

import (
	"errors"
	"fmt"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	ErrUnknownOpcode = "unknown opcode: %s"
)

// Closure encapsulates a function with the locals it captures
type Closure struct {
	*Lambda
	data.Call
	data.Values
}

func newClosure(lambda *Lambda, values data.Values) *Closure {
	c := &Closure{
		Lambda: lambda,
		Values: values,
	}

	c.Call = func(args ...data.Value) data.Value {
		closure := c
		lambda := closure.Lambda
		code := lambda.Code
		stackSize := lambda.StackSize
		localCount := lambda.LocalCount
		stack := make(data.Values, stackSize)
		locals := make(data.Values, localCount)
		var SP = stackSize - 1
		var PC = 0
		goto opSwitch

	nextPC:
		PC++

	opSwitch:
		op := isa.Opcode(code[PC])
		switch op {
		case isa.Null:
			stack[SP] = data.Null
			SP--
			goto nextPC

		case isa.Zero:
			stack[SP] = data.Integer(0)
			SP--
			goto nextPC

		case isa.One:
			stack[SP] = data.Integer(1)
			SP--
			goto nextPC

		case isa.NegOne:
			stack[SP] = data.Integer(-1)
			SP--
			goto nextPC

		case isa.Two:
			stack[SP] = data.Integer(2)
			SP--
			goto nextPC

		case isa.True:
			stack[SP] = data.True
			SP--
			goto nextPC

		case isa.False:
			stack[SP] = data.False
			SP--
			goto nextPC

		case isa.Const:
			PC++
			idx := isa.Index(code[PC])
			stack[SP] = lambda.Constants[idx]
			SP--
			goto nextPC

		case isa.Arg:
			PC++
			idx := isa.Index(code[PC])
			stack[SP] = args[idx]
			SP--
			goto nextPC

		case isa.RestArg:
			PC++
			idx := isa.Index(code[PC])
			stack[SP] = data.NewVector(args[idx:]...)
			SP--
			goto nextPC

		case isa.ArgLen:
			stack[SP] = data.Integer(len(args))
			SP--
			goto nextPC

		case isa.Closure:
			PC++
			idx := isa.Index(code[PC])
			stack[SP] = closure.Values[idx]
			SP--
			goto nextPC

		case isa.Load:
			PC++
			idx := isa.Index(code[PC])
			stack[SP] = locals[idx]
			SP--
			goto nextPC

		case isa.Store:
			PC++
			idx := isa.Index(code[PC])
			SP++
			locals[idx] = stack[SP]
			goto nextPC

		case isa.NewRef:
			stack[SP] = &Ref{}
			SP--
			goto nextPC

		case isa.BindRef:
			SP++
			ref := stack[SP].(*Ref)
			SP++
			ref.Value = stack[SP].(data.Value)
			goto nextPC

		case isa.Deref:
			SP1 := SP + 1
			stack[SP1] = stack[SP1].(*Ref).Value
			goto nextPC

		case isa.Resolve:
			SP1 := SP + 1
			sym := stack[SP1].(data.Symbol)
			val := namespace.MustResolveValue(lambda.Globals, sym)
			stack[SP1] = val
			goto nextPC

		case isa.Declare:
			SP++
			name := stack[SP].(data.Name)
			lambda.Globals.Declare(name)
			goto nextPC

		case isa.Bind:
			SP++
			name := stack[SP].(data.Name)
			SP++
			val := stack[SP].(data.Value)
			lambda.Globals.Declare(name).Bind(val)
			goto nextPC

		case isa.Dup:
			stack[SP] = stack[SP+1]
			SP--
			goto nextPC

		case isa.Pop:
			SP++
			goto nextPC

		case isa.Add:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			stack[SP1] = left.Add(right)
			goto nextPC

		case isa.Sub:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			stack[SP1] = left.Sub(right)
			goto nextPC

		case isa.Mul:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			stack[SP1] = left.Mul(right)
			goto nextPC

		case isa.Div:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			stack[SP1] = left.Div(right)
			goto nextPC

		case isa.Mod:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			stack[SP1] = left.Mod(right)
			goto nextPC

		case isa.Eq:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			cmp := left.Cmp(right)
			stack[SP1] = data.Bool(cmp == data.EqualTo)
			goto nextPC

		case isa.Neq:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			cmp := left.Cmp(right)
			stack[SP1] = data.Bool(cmp != data.EqualTo)
			goto nextPC

		case isa.Lt:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			cmp := left.Cmp(right)
			stack[SP1] = data.Bool(cmp == data.LessThan)
			goto nextPC

		case isa.Lte:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			cmp := left.Cmp(right)
			res := cmp == data.LessThan || cmp == data.EqualTo
			stack[SP1] = data.Bool(res)
			goto nextPC

		case isa.Gt:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			cmp := left.Cmp(right)
			stack[SP1] = data.Bool(cmp == data.GreaterThan)
			goto nextPC

		case isa.Gte:
			SP++
			SP1 := SP + 1
			right := stack[SP].(data.Number)
			left := stack[SP1].(data.Number)
			cmp := left.Cmp(right)
			res := cmp == data.GreaterThan || cmp == data.EqualTo
			stack[SP1] = data.Bool(res)
			goto nextPC

		case isa.Neg:
			SP1 := SP + 1
			val := stack[SP1].(data.Number)
			stack[SP1] = data.Integer(0).Sub(val)
			goto nextPC

		case isa.Not:
			SP1 := SP + 1
			val := stack[SP1].(data.Bool)
			stack[SP1] = data.Bool(!val)
			goto nextPC

		case isa.MakeTruthy:
			SP1 := SP + 1
			val := data.Truthy(stack[SP1])
			stack[SP1] = data.Bool(val)
			goto nextPC

		case isa.MakeCall:
			SP1 := SP + 1
			val := stack[SP1]
			if _, ok := val.(data.Call); ok {
				goto nextPC
			}
			stack[SP1] = val.(data.Caller).Caller()
			goto nextPC

		case isa.Call0:
			SP1 := SP + 1
			fn := stack[SP1].(data.Call)
			stack[SP1] = fn()
			goto nextPC

		case isa.Call1:
			SP++
			SP1 := SP + 1
			fn := stack[SP].(data.Call)
			arg := stack[SP1]
			stack[SP1] = fn(arg)
			goto nextPC

		case isa.Call:
			PC++
			SP1 := SP + 1
			SP2 := SP1 + 1
			fn := stack[SP1].(data.Call)
			argCount := isa.Count(code[PC])
			RES := SP1 + int(argCount)
			args := make(data.Values, argCount)
			copy(args, stack[SP2:])
			stack[RES] = fn(args...)
			SP = RES - 1
			goto nextPC

		case isa.TailCall:
			SP1 := SP + 1
			argCount := int(code[PC+1])
			args = make(data.Values, argCount)
			copy(args, stack[SP1+1:])
			val := stack[SP1]
			if vc, ok := val.(*Closure); ok {
				if vc != closure {
					closure = vc
					lambda = closure.Lambda
					code = lambda.Code
					stackSize = lambda.StackSize
					localCount = lambda.LocalCount
					stack = make(data.Values, stackSize)
					locals = make(data.Values, localCount)
				}
				SP = stackSize - 1
				PC = 0
				goto opSwitch
			}
			fn := val.(data.Caller)
			return fn.Caller()(args...)

		case isa.Jump:
			off := isa.Offset(code[PC+1])
			PC = int(off)
			goto opSwitch

		case isa.CondJump:
			SP++
			val := stack[SP].(data.Bool)
			if val {
				off := isa.Offset(code[PC+1])
				PC = int(off)
				goto opSwitch
			}
			PC += 2
			goto opSwitch

		case isa.Panic:
			panic(errors.New(stack[SP+1].String()))

		case isa.Return:
			return stack[SP+1]

		case isa.RetNull:
			return data.Null

		case isa.RetEmptyList:
			return data.EmptyList

		case isa.RetTrue:
			return data.True

		case isa.RetFalse:
			return data.False

		default:
			panic(fmt.Errorf(ErrUnknownOpcode, op))
		}
	}

	return c
}

// Caller returns a calling interface for this Closure
func (c *Closure) Caller() data.Call {
	return c.Call
}

// CheckArity performs a compile-time arity check for the closure
func (c *Closure) CheckArity(i int) error {
	return c.ArityChecker(i)
}

// Convention returns the closure's calling convention
func (c *Closure) Convention() data.Convention {
	return data.ApplicativeCall
}

// Type makes Closure a typed value
func (c *Closure) Type() data.Name {
	res := fmt.Sprintf("%s-closure", c.Convention())
	return data.Name(res)
}

func (c *Closure) String() string {
	return data.DumpString(c)
}
