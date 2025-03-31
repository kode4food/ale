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
	Add      Opcode = iota // Addition
	Arg                    // Retrieve Argument Value
	ArgLen                 // Retrieve Argument Count
	Bind                   // Bind Global
	BindRef                // Bind Reference
	Call                   // Call Procedure
	Call0                  // Zero-Arg Call
	Call1                  // One-Arg Call
	Call2                  // Two-Arg Call
	Call3                  // Three-Arg Call
	CallWith               // Call with the provided Sequence as Arguments
	Car                    // Contents of the Address part of the Register
	Cdr                    // Contents of the Decrement part of the Register
	Closure                // Retrieve Closure Value
	CondJump               // Conditional Jump
	Cons                   // Form a new Cons Cell (or Prepend)
	Const                  // Retrieve Constant
	Deref                  // Pointer Dereference
	Div                    // Division
	Dup                    // Duplicate Value
	Empty                  // Tests Empty Sequence
	Eq                     // Value Equality
	False                  // Push False
	Jump                   // Absolute Jump
	Label                  // Internal Label
	Load                   // Retrieve Local Value
	Mod                    // Remainder
	Mul                    // Multiplication
	Neg                    // Negation
	NegInt                 // Push Negative Integer (in Operand)
	NewRef                 // New Reference
	NoOp                   // Error-Operator
	Not                    // Boolean Negation
	Null                   // Push Null
	NumEq                  // Numeric Equality
	NumGt                  // Greater Than
	NumGte                 // Greater or Equal To
	NumLt                  // Less Than Comparison
	NumLte                 // Less or Equal To
	Panic                  // Halt Abnormally
	Pop                    // Discard Value
	PopArgs                // Pop Arguments
	PosInt                 // Push Positive Integer (in Operand)
	Private                // Declare a private Namespace entry
	Public                 // Declare a public Namespace entry
	PushArgs               // Push Arguments
	Resolve                // Resolve Global Symbol
	RestArg                // Retrieve Arguments Tail
	RetFalse               // Return False
	RetNull                // Return Null (Empty List)
	RetTrue                // Return True
	Return                 // Return Value
	Store                  // Store Local
	Sub                    // Subtraction
	TailCall               // Tail Call
	True                   // Push True
	Vector                 // Make a new Vector from the Stack
	Zero                   // Push Zero
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
