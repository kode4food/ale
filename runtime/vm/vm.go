package vm

import (
	"errors"
	"fmt"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

type (
	// Config encapsulates the initial environment of a virtual machine
	Config struct {
		Globals    namespace.Type
		Constants  data.Values
		Code       []isa.Word
		StackSize  int
		LocalCount int
	}

	// Closure passes enclosed state into a Caller
	Closure func(...data.Value) data.Call
)

// Error messages
const (
	ErrUnknownOpcode = "unknown opcode: %s"
)

// NewClosure returns a Closure based on the virtual machine configuration
func NewClosure(cfg *Config) data.Call {
	globals := cfg.Globals
	constants := cfg.Constants
	code := cfg.Code
	stackSize := cfg.StackSize
	localCount := cfg.LocalCount
	stackInit := stackSize - 1

	return func(closure ...data.Value) data.Value {
		var self data.Call

		self = func(args ...data.Value) data.Value {
			stack := make(data.Values, stackSize)
			locals := make(data.Values, localCount)
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
				stack[SP] = data.Nil
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
				stack[SP] = data.Vector(args[idx:])
				SP--
				goto nextPC

			case isa.ArgLen:
				stack[SP] = data.Integer(len(args))
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
				sym := stack[SP1].(data.Symbol)
				val := namespace.MustResolveValue(globals, sym)
				stack[SP1] = val
				goto nextPC

			case isa.Declare:
				SP++
				name := stack[SP].(data.Name)
				globals.Declare(name)
				goto nextPC

			case isa.Bind:
				SP++
				name := stack[SP].(data.Name)
				SP++
				val := stack[SP].(data.Value)
				globals.Declare(name).Bind(val)
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
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Integer(left + right)
				goto nextPC

			case isa.Sub:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Integer(left - right)
				goto nextPC

			case isa.Mul:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Integer(left * right)
				goto nextPC

			case isa.Div:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Integer(left / right)
				goto nextPC

			case isa.Mod:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Integer(left % right)
				goto nextPC

			case isa.Eq:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Bool(left == right)
				goto nextPC

			case isa.Neq:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Bool(left != right)
				goto nextPC

			case isa.Lt:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Bool(left < right)
				goto nextPC

			case isa.Lte:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Bool(left <= right)
				goto nextPC

			case isa.Gt:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Bool(left > right)
				goto nextPC

			case isa.Gte:
				SP++
				SP1 := SP + 1
				right := stack[SP].(data.Integer)
				left := stack[SP1].(data.Integer)
				stack[SP1] = data.Bool(left >= right)
				goto nextPC

			case isa.Neg:
				SP1 := SP + 1
				val := stack[SP1].(data.Integer)
				stack[SP1] = data.Integer(-val)
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
				val := stack[SP1].(data.Caller)
				stack[SP1] = val.Caller()
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
				args := make([]data.Value, argCount)
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

			case isa.RetNil:
				return data.Nil

			case isa.RetTrue:
				return data.True

			case isa.RetFalse:
				return data.False

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
