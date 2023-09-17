package isa

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

	// Label is an internal Opcode
	Label = OpcodeMask
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=Opcode
const (
	Add        Opcode = iota // Addition
	Arg                      // Retrieve Argument Value
	ArgLen                   // Retrieve Argument Count
	Bind                     // Bind Global
	BindRef                  // Bind Reference
	Call                     // Call Lambda
	Call0                    // Zero-Arg Call
	Call1                    // One-Arg Call
	Car                      // Contents of the Address part of the Register
	Cdr                      // Contents of the Decrement part of the Register
	Closure                  // Retrieve Closure Value
	CondJump                 // Conditional Jump
	Const                    // Retrieve Constant
	Declare                  // Declare a public Namespace entry
	Deref                    // Pointer Dereference
	Div                      // Division
	Dup                      // Duplicate Value
	Eq                       // Numeric Equality
	False                    // Push False
	Gt                       // Greater Than
	Gte                      // Greater or Equal To
	Jump                     // Absolute Jump
	Load                     // Retrieve Local Value
	Lt                       // Less Than Comparison
	Lte                      // Less or Equal To
	MakeTruthy               // Make Value Boolean
	Mod                      // Remainder
	Mul                      // Multiplication
	Neg                      // Negation
	NegInf                   // Push Negative Infinity
	NegInt                   // Push Negative Integer (in Operand)
	Neq                      // Numeric Inequality
	NewRef                   // New Reference
	NoOp                     // Non-Operator
	Not                      // Boolean Negation
	Nil                      // Push Nil
	Panic                    // Abnormally Halt
	Pop                      // Discard Value
	PosInf                   // Positive Infinity
	PosInt                   // Push Positive Integer (in Operand)
	Private                  // Declare a private Namespace entry
	Resolve                  // Resolve Global Symbol
	RestArg                  // Retrieve Arguments Tail
	RetFalse                 // Return False
	RetNil                   // Return Nil (Empty List)
	RetTrue                  // Return True
	Return                   // Return Value
	Self                     // Push Current Lambda
	Store                    // Store Local
	Sub                      // Subtraction
	TailCall                 // Tail Call
	True                     // Push True
	Zero                     // Push Zero
)

func (o Opcode) New(ops ...Operand) Instruction {
	return New(o, ops...)
}
