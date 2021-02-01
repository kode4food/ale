package isa

// Opcode represents an Instruction's operation Word
type Opcode Word

// Label is an internal Opcode
const Label Opcode = 256

//go:generate stringer -type=Opcode
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
	Declare                  // Declare Global
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
	NegOne                   // Push Negative One
	Neq                      // Numeric Inequality
	NewRef                   // New Reference
	NoOp                     // Non-Operator
	Not                      // Boolean Negation
	Nil                      // Push Nil
	One                      // Push One
	Panic                    // Abnormally Halt
	Pop                      // Discard Value
	PosInf                   // Positive Infinity
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
	Two                      // Push Two
	Zero                     // Push Zero
)

// Word makes Opcode a Coder
func (i Opcode) Word() Word {
	return Word(i)
}
