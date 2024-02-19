package vm

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	State int

	VM struct {
		CL   *Closure
		ST   State
		PC   int
		LP   int
		SP   int
		INST isa.Instruction
		CODE isa.Instructions
		MEM  data.Vector
		ARGS ArgStack
		RES  data.Value
	}

	ArgStack struct {
		Data data.Vector
		Next *ArgStack
	}
)

// ErrBadInstruction is raised when the VM encounters an Opcode that has not
// been properly mapped
const ErrBadInstruction = "unknown instruction encountered: %s"

const (
	FAILURE State = iota - 1
	RUNNING
	SUCCESS
)

func (vm *VM) initMem() {
	vm.MEM = make(data.Vector, vm.CL.StackSize+vm.CL.LocalCount)
	vm.initCode()
}

func (vm *VM) initCode() {
	vm.CODE = vm.CL.Code
	vm.LP = int(vm.CL.StackSize)
	vm.initState()
}

func (vm *VM) initState() {
	vm.SP = vm.LP - 1
	vm.PC = 0
}

func (vm *VM) Run() data.Value {
	vm.initMem()
	for vm.ST == RUNNING {
		vm.INST = vm.CODE[vm.PC]
		switch vm.INST.Opcode() {
		case isa.Add:
			doAdd(vm)
		case isa.Arg:
			doArg(vm)
		case isa.ArgLen:
			doArgLen(vm)
		case isa.Bind:
			doBind(vm)
		case isa.BindRef:
			doBindRef(vm)
		case isa.Call0:
			doCall0(vm)
		case isa.Call1:
			doCall1(vm)
		case isa.Call:
			doCall(vm)
		case isa.CallWith:
			doCallWith(vm)
		case isa.Car:
			doCar(vm)
		case isa.Cdr:
			doCdr(vm)
		case isa.Closure:
			doClosure(vm)
		case isa.CondJump:
			doCondJump(vm)
		case isa.Cons:
			doCons(vm)
		case isa.Const:
			doConst(vm)
		case isa.Declare:
			doDeclare(vm)
		case isa.Deref:
			doDeref(vm)
		case isa.Div:
			doDiv(vm)
		case isa.Dup:
			doDup(vm)
		case isa.Empty:
			doEmpty(vm)
		case isa.Eq:
			doEq(vm)
		case isa.False:
			doFalse(vm)
		case isa.Jump:
			doJump(vm)
		case isa.Load:
			doLoad(vm)
		case isa.Mod:
			doMod(vm)
		case isa.Mul:
			doMul(vm)
		case isa.Neg:
			doNeg(vm)
		case isa.NegInt:
			doNegInt(vm)
		case isa.NewRef:
			doNewRef(vm)
		case isa.NoOp:
			doNoOp(vm)
		case isa.Not:
			doNot(vm)
		case isa.Null:
			doNull(vm)
		case isa.NumEq:
			doNumEq(vm)
		case isa.NumGt:
			doNumGt(vm)
		case isa.NumGte:
			doNumGte(vm)
		case isa.NumLt:
			doNumLt(vm)
		case isa.NumLte:
			doNumLte(vm)
		case isa.Panic:
			doPanic(vm)
		case isa.Pop:
			doPop(vm)
		case isa.PopArgs:
			doPopArgs(vm)
		case isa.PosInt:
			doPosInt(vm)
		case isa.Private:
			doPrivate(vm)
		case isa.PushArgs:
			doPushArgs(vm)
		case isa.Resolve:
			doResolve(vm)
		case isa.RestArg:
			doRestArg(vm)
		case isa.RetFalse:
			doRetFalse(vm)
		case isa.RetNull:
			doRetNull(vm)
		case isa.RetTrue:
			doRetTrue(vm)
		case isa.Return:
			doReturn(vm)
		case isa.Store:
			doStore(vm)
		case isa.Sub:
			doSub(vm)
		case isa.TailCall:
			doTailCall(vm)
		case isa.True:
			doTrue(vm)
		case isa.Zero:
			doZero(vm)
		default:
			panic(debug.ProgrammerError(ErrBadInstruction, vm.INST))
		}
	}
	return vm.RES
}
