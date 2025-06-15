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
	OpcodeSize = 6

	// OpcodeMask masks the bits for encoding an Opcode into an Instruction
	OpcodeMask = Opcode(1<<OpcodeSize - 1)

	// OperandMask masks the bits for encoding an Operand into an Instruction
	OperandMask = ^Operand(OpcodeMask) >> OpcodeSize
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=Opcode
const (
	Add        Opcode = iota // Pop two numbers, push their sum
	Arg                      // Push the Nth argument (op = index)
	ArgsLen                  // Push the number of arguments
	ArgsPop                  // Pop argument stack frame
	ArgsPush                 // Push N arguments (op = count)
	Call                     // Calls proc with N args (op = count)
	Call0                    // Calls proc with zero arguments
	Call1                    // Calls proc with one argument
	Call2                    // Calls proc with two arguments
	Call3                    // Calls proc with three arguments
	CallSelf                 // Calls current closure with N args (op = count)
	CallWith                 // Calls proc with args from a sequence
	Car                      // Pop pair, push its Address (car) part
	Cdr                      // Pop pair, push its Decrement (cdr) part
	Closure                  // Push captured value (op = index)
	CondJump                 // If top of stack is true, jumps to operand
	Cons                     // Pop two values, push a new cons cell
	Const                    // Push constant (op = index)
	Div                      // Pop two numbers, push their quotient
	Dup                      // Duplicates the top of the stack
	Empty                    // Push true if sequence at top of stack is empty
	Eq                       // Pop two values, push true if equal
	False                    // Push the boolean false
	EnvBind                  // Pop name and value, binds value to global symbol
	EnvPrivate               // Pop symbol, marks as private in namespace
	EnvPublic                // Pop symbol, marks as public in namespace
	EnvValue                 // Resolves namespace symbol at top of stack
	Jump                     // Jumps to operand
	Label                    // Marks a label (not executed)
	Load                     // Push local value (op = index)
	Mod                      // Pop two numbers, push remainder of division
	Mul                      // Pop two numbers, push their product
	Neg                      // Negates the number at top of stack
	NegInt                   // Push negative integer (operand)
	NewRef                   // Push a new unbound Ref
	NoOp                     // No operation
	Not                      // Boolean negation of top of stack
	Null                     // Push the null value
	NumEq                    // Pop two numbers, push true if equal
	NumGt                    // Pop two numbers, push true if second > first
	NumGte                   // Pop two numbers, push true if second >= first
	NumLt                    // Pop two numbers, push true if second < first
	NumLte                   // Pop two numbers, push true if second <= first
	Panic                    // Raises error with value at top of stack
	Pop                      // Discards the top of the stack
	PosInt                   // Push positive integer (operand)
	RefBind                  // Pop ref and value, sets ref.Value to value
	RefValue                 // Loads ref.Value from ref at top of stack
	RestArg                  // Push rest of arguments as a vector
	RetFalse                 // Returns the boolean false
	RetNull                  // Returns the null value
	RetTrue                  // Returns the boolean true
	Return                   // Returns value at top of stack
	Store                    // Pop value, stores in local (op = index)
	Sub                      // Pop two numbers, push their difference
	TailCall                 // Dynamic tail call with N args (op = count)
	TailClos                 // Tail call to a closure with N args (op = count)
	TailSelf                 // Tail call to self with N args (op = count)
	True                     // Push the boolean true
	Vector                   // Pop N values, push as a vector (op = count)
	Zero                     // Push the integer zero
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
