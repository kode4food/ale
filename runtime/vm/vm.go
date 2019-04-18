package vm

import (
	"errors"
	"fmt"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

type (
	// Config encapsulates the initial environment of a virtual machine
	Config struct {
		Globals    namespace.Type
		Constants  api.Values
		Code       []isa.Word
		StackSize  int
		LocalCount int
	}

	// Closure passes enclosed state into a Caller
	Closure func(...api.Value) api.Call
)

// Error messages
const (
	ErrUnknownOpcode = "unknown opcode: %s"
)

// NewClosure returns a Closure based on the virtual machine configuration
func NewClosure(cfg *Config) api.Call {
	globals := cfg.Globals
	constants := cfg.Constants
	code := cfg.Code
	stackSize := cfg.StackSize
	localCount := cfg.LocalCount
	stackInit := stackSize - 1

	return func(closure ...api.Value) api.Value {
		var self api.Call

		self = func(args ...api.Value) api.Value {
			stack := make(api.Values, stackSize)
			locals := make(api.Values, localCount)
			var PC = 0
			var SP = stackInit
			goto opSwitch

		nextPC:
			PC++

		opSwitch:
			op := isa.Opcode(code[PC])
			switch op {
			case isa.NoOp:
				goto nextPC

			case isa.Self:
				stack[SP] = self
				SP--
				goto nextPC

			case isa.Nil:
				stack[SP] = api.Nil
				SP--
				goto nextPC

			case isa.Zero:
				stack[SP] = api.Integer(0)
				SP--
				goto nextPC

			case isa.One:
				stack[SP] = api.Integer(1)
				SP--
				goto nextPC

			case isa.NegOne:
				stack[SP] = api.Integer(-1)
				SP--
				goto nextPC

			case isa.Two:
				stack[SP] = api.Integer(2)
				SP--
				goto nextPC

			case isa.True:
				stack[SP] = api.True
				SP--
				goto nextPC

			case isa.False:
				stack[SP] = api.False
				SP--
				goto nextPC

			case isa.Const:
				PC++
				idx := isa.Index(code[PC])
				stack[SP] = constants[idx]
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
				stack[SP] = api.Vector(args[idx:])
				SP--
				goto nextPC

			case isa.ArgLen:
				stack[SP] = api.Integer(len(args))
				SP--
				goto nextPC

			case isa.Closure:
				PC++
				idx := isa.Index(code[PC])
				stack[SP] = closure[idx]
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

			case isa.Resolve:
				SP1 := SP + 1
				sym := stack[SP1].(api.Symbol)
				val := namespace.MustResolveSymbol(globals, sym)
				stack[SP1] = val
				goto nextPC

			case isa.Declare:
				SP++
				name := stack[SP].(api.Name)
				globals.Declare(name)
				goto nextPC

			case isa.Bind:
				SP++
				name := stack[SP].(api.Name)
				SP++
				val := stack[SP].(api.Value)
				globals.Bind(name, val)
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
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Integer(left + right)
				goto nextPC

			case isa.Sub:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Integer(left - right)
				goto nextPC

			case isa.Mul:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Integer(left * right)
				goto nextPC

			case isa.Div:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Integer(left / right)
				goto nextPC

			case isa.Mod:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Integer(left % right)
				goto nextPC

			case isa.Eq:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Bool(left == right)
				goto nextPC

			case isa.Neq:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Bool(left != right)
				goto nextPC

			case isa.Lt:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Bool(left < right)
				goto nextPC

			case isa.Lte:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Bool(left <= right)
				goto nextPC

			case isa.Gt:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Bool(left > right)
				goto nextPC

			case isa.Gte:
				SP++
				SP1 := SP + 1
				right := stack[SP].(api.Integer)
				left := stack[SP1].(api.Integer)
				stack[SP1] = api.Bool(left >= right)
				goto nextPC

			case isa.Neg:
				SP1 := SP + 1
				val := stack[SP1].(api.Integer)
				stack[SP1] = api.Integer(-val)
				goto nextPC

			case isa.Not:
				SP1 := SP + 1
				val := stack[SP1].(api.Bool)
				stack[SP1] = api.Bool(!val)
				goto nextPC

			case isa.MakeTruthy:
				SP1 := SP + 1
				val := api.Truthy(stack[SP1])
				stack[SP1] = api.Bool(val)
				goto nextPC

			case isa.MakeCall:
				SP1 := SP + 1
				val := stack[SP1].(api.Caller)
				stack[SP1] = val.Caller()
				goto nextPC

			case isa.Call0:
				SP1 := SP + 1
				fn := stack[SP1].(api.Call)
				stack[SP1] = fn()
				goto nextPC

			case isa.Call1:
				SP++
				SP1 := SP + 1
				fn := stack[SP].(api.Call)
				arg := stack[SP1]
				stack[SP1] = fn(arg)
				goto nextPC

			case isa.Call:
				PC++
				SP1 := SP + 1
				SP2 := SP1 + 1
				fn := stack[SP1].(api.Call)
				argCount := isa.Count(code[PC])
				RES := SP1 + int(argCount)
				args := make([]api.Value, argCount)
				copy(args, stack[SP2:])
				stack[RES] = fn(args...)
				SP = RES - 1
				goto nextPC

			case isa.Jump:
				off := isa.Offset(code[PC+1])
				PC = int(off)
				goto opSwitch

			case isa.CondJump:
				SP++
				val := stack[SP].(api.Bool)
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

			case isa.RetNil:
				return api.Nil

			case isa.RetTrue:
				return api.True

			case isa.RetFalse:
				return api.False

			default:
				panic(fmt.Errorf(ErrUnknownOpcode, op))
			}
		}

		return self
	}
}

func (Closure) String() string {
	return "closure"
}
