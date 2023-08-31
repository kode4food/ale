package isa

import "fmt"

// Opcode represents an Instruction's operation
type Opcode Word

const (
	OpcodeMask  = 0x3F      // 6 bit mask
	OpcodeSize  = 6         // number of bits to shift
	OperandMask = 0x3FFFFFF // 26 bit mask

	// Label is an internal Opcode
	Label Opcode = Opcode(OpcodeMask)
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

func (oc Opcode) Instruction() Instruction {
	if f, ok := Effects[oc]; ok && f.Operand == Nothing {
		return Instruction(oc)
	}
	// Programmer error
	panic(fmt.Sprintf("opcode can't be encoded as instruction: %s", oc))
}
