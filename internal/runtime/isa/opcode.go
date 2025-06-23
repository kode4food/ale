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

	// Argument, Environemnt and Closure Operations
	Arg        // Push the Nth argument (op = index)
	ArgsLen    // Push the number of arguments
	ArgsPop    // Pop argument stack frame
	ArgsPush   // Push N arguments (op = count)
	ArgsRest   // Push rest of arguments as a vector
	Closure    // Push captured value (op = index)
	EnvBind    // Pop name and value, binds value to global symbol
	EnvPrivate // Pop symbol, marks as private in namespace
	EnvPublic  // Pop symbol, marks as public in namespace
	EnvValue   // Resolves namespace symbol at top of stack

	// Reference and Register Operations
	Load     // Push local value (op = index)
	NewRef   // Push a new unbound Ref
	RefBind  // Pop ref and value, sets ref.Value to value
	RefValue // Loads ref.Value from ref at top of stack
	Store    // Pop value, stores in local (op = index)

	// Stack and Constant Operations
	Dup   // Duplicates the top of the stack
	Pop   // Discards the top of the stack
	Const // Push constant (op = index)
	False // Push the boolean false
	Null  // Push the null value
	True  // Push the boolean true
	Zero  // Push the integer zero

	// Call Operations
	Call     // Calls proc with N args (op = count)
	Call0    // Calls proc with zero arguments
	Call1    // Calls proc with one argument
	Call2    // Calls proc with two arguments
	Call3    // Calls proc with three arguments
	CallSelf // Calls current closure with N args (op = count)
	CallWith // Calls proc with args from a sequence
	TailCall // Dynamic tail call with N args (op = count)
	TailClos // Tail call to a closure with N args (op = count)
	TailSelf // Tail call to self with N args (op = count)

	// Control Flow Operations
	CondJump // If top of stack is true, jumps to operand
	Jump     // Jumps to operand
	Panic    // Raises error with value at top of stack
	RetFalse // Returns the boolean false
	RetNull  // Returns the null value
	RetTrue  // Returns the boolean true
	Return   // Returns value at top of stack

	// Sequence Operations
	Append // Pop value, append to sequence at top of stack
	Assoc  // Pop pair, associate with mapper at top of stack
	Car    // Pop pair, push its Address (car) part
	Cdr    // Pop pair, push its Decrement (cdr) part
	Cons   // Pop two values, push a new cons cell
	Dissoc // Pop key, dissociate from mapper at top of stack
	Empty  // Push true if sequence at top of stack is empty
	Get    // Pop key, push value from mapper at top of stack
	Nth    // Pop index, push value from sequence at top of stack
	Vector // Pop N values, push as a vector (op = count)

	// Boolean Operations
	Eq  // Pop two values, push true if equal
	Not // Boolean negation of top of stack

	// Numeric Operations
	Add    // Pop two numbers, push their sum
	Div    // Pop two numbers, push their quotient
	Mod    // Pop two numbers, push remainder of division
	Mul    // Pop two numbers, push their product
	Neg    // Negates the number at top of stack
	NegInt // Push negative integer (operand)
	NumEq  // Pop two numbers, push true if equal
	NumGt  // Pop two numbers, push true if second > first
	NumGte // Pop two numbers, push true if second >= first
	NumLt  // Pop two numbers, push true if second < first
	NumLte // Pop two numbers, push true if second <= first
	PosInt // Push positive integer (operand)
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
