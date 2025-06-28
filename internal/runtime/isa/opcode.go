package isa

import "fmt"

type (
	// Opcode represents an Instruction's operation
	Opcode Word

	// Operand allows an Instruction to be parameterized
	Operand Word
)

const (
	// OpcodeSize are the number of bits required for an Opcode value
	OpcodeSize = 7

	// OpcodeMask masks the bits for encoding an Opcode into an Instruction
	OpcodeMask = Opcode(1<<OpcodeSize - 1)

	// OperandMask masks the bits for encoding an Operand into an Instruction
	OperandMask = ^Operand(OpcodeMask) >> OpcodeSize
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=Opcode
const (
	// Ignored Opcodes
	Label Opcode = iota // Marks a label (not executed)
	NoOp                // No operation

	// Argument, Environment and Closure Operations
	Arg        // Push the Nth argument (op = index)
	ArgsLen    // Push the number of arguments
	ArgsPop    // Pop argument stack, restore previous arguments
	ArgsPush   // Push argument stack, replace with popped values (op = count)
	ArgsRest   // Push remaining arguments as a vector (op = start index)
	Closure    // Push captured value (op = index)
	EnvBind    // Pop symbol and value, binds value to namespace symbol
	EnvPrivate // Pop symbol, mark as private in namespace
	EnvPublic  // Pop symbol, mark as public in namespace
	EnvValue   // Pop symbol, resolve it from namespace

	// Reference and Register Operations
	Load     // Push local value (op = index)
	NewRef   // Push a new unbound Ref
	RefBind  // Pop ref and value, sets ref.Value to value
	RefValue // Pop ref, push ref.Value
	Store    // Pop value, store in local (op = index)

	// Stack and Constant Operations
	Const // Push constant (op = index)
	Dup   // Push a duplicate of the top of the stack
	False // Push the boolean false
	Null  // Push the null value
	Pop   // Pop (discard) the top of the stack
	Swap  // Swap the top two values on the stack
	True  // Push the boolean true
	Zero  // Push the integer zero

	// Call Operations
	Call     // Pop proc, pop N args, call proc, push result (op = count)
	Call0    // Pop proc, call with zero args, push result
	Call1    // Pop proc, pop 1 arg, call proc, push result
	Call2    // Pop proc, pop 2 args, call proc, push result
	Call3    // Pop proc, pop 3 args, call proc, push result
	CallSelf // Pop N args, call current closure (op = count)
	CallWith // Pop proc, pop sequence, call proc with seq values (op = count)
	TailCall // Pop proc, pop N args, dynamic tail call (op = count)
	TailClos // Pop closure, pop N args, tail call (op = count)
	TailSelf // Pop N args, tail call (op = count)

	// Control Flow Operations
	CondJump // Pop value, if not false, jump to operand
	Delay    // Pop proc, wrap as promise, push promise
	Jump     // Jump to operand
	Panic    // Pop value, raise as error from VM
	RetFalse // Return the boolean false from VM
	RetNull  // Return the null value from VM
	RetTrue  // Return the boolean true from VM
	Return   // Pop value, return value from VM

	// Sequence Operations
	Append  // Pop value, pop sequence, append value to sequence, push result
	Assoc   // Pop pair, pop mapper, associate pair with mapper, push result
	Car     // Pop pair, push pair's Address (car) part
	Cdr     // Pop pair, push pair's Decrement (cdr) part
	Cons    // Pop car, pop cdr, push new cons cell
	Dissoc  // Pop key, pop mapper, dissociate key from mapper, push result
	Empty   // Pop sequence, push true if sequence is empty
	Get     // Pop key, pop mapper, push found status and value from mapper
	LazySeq // Pop proc, wrap as lazy sequence, push lazy sequence
	Length  // Pop sequence, push sequence length
	Nth     // Pop index, pop indexed, push from status and value from indexed
	Reverse // Pop sequence, push reversed sequence
	Vector  // Pop N values, push as a vector (op = count)

	// Boolean Operations
	Eq  // Pop two values, push true if equal
	Not // Pop value, push boolean negation of value

	// Numeric Operations
	Add    // Pop two numbers, push their sum
	Div    // Pop two numbers, push their quotient
	Mod    // Pop two numbers, push the remainder of division
	Mul    // Pop two numbers, push their product
	Neg    // Pop number, push negated value of number
	NegInt // Push negative integer (int = operand)
	NumEq  // Pop two numbers, push true if equal
	NumGt  // Pop two numbers, push true if second > first
	NumGte // Pop two numbers, push true if second >= first
	NumLt  // Pop two numbers, push true if second < first
	NumLte // Pop two numbers, push true if second <= first
	PosInt // Push positive integer (int = operand)
	Sub    // Pop two numbers, push their difference
)

// New creates a new Instruction instance from an Opcode
func (o Opcode) New(ops ...Operand) Instruction {
	res, err := o.new(ops...)
	if err != nil {
		panic(err)
	}
	return res
}

func (o Opcode) new(ops ...Operand) (Instruction, error) {
	effect, err := GetEffect(o)
	if err != nil {
		return 0, err
	}
	switch {
	case effect.Operand != Nothing && len(ops) == 1:
		if !IsValidOperand(int(ops[0])) {
			return 0, fmt.Errorf(ErrExpectedOperand, ops[0])
		}
		return Instruction(Opcode(ops[0]<<OpcodeSize) | o), nil
	case effect.Operand == Nothing && len(ops) == 0:
		return Instruction(o), nil
	default:
		return 0, fmt.Errorf(ErrBadInstruction, o.String())
	}
}
