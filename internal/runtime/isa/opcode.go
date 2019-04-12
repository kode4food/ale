package isa

// Opcode represents an Instruction's operation Word
type Opcode Word

// Label is an internal Opcode
const Label Opcode = 256

//go:generate stringer -type=Opcode -linecomment
const (
	Add         Opcode = iota // Addition
	Arg                       // Retrieve Argument Value
	ArgLen                    // Retrieve Argument Count
	Bind                      // Bind Global
	Call                      // Call Function
	Call0                     // Zero-Arg Call
	Call1                     // One-Arg Call
	Closure                   // Retrieve Closure Value
	CondJump                  // Conditional Jump
	Const                     // Retrieve Constant
	Declare                   // Declare Global
	Div                       // Division
	Dup                       // Duplicate Value
	Eq                        // Numeric Equality
	False                     // Push False
	Gt                        // Greater Than
	Gte                       // Greater or Equal To
	Jump                      // Absolute Jump
	Load                      // Retrieve Local Value
	Lt                        // Less Than Comparison
	Lte                       // Less or Equal To
	MakeCall                  // Make Value Callable
	MakeTruthy                // Make Value Boolean
	Mod                       // Remainder
	Mul                       // Multiplication
	Neg                       // Negation
	NegInfinity               // Push Negative Infinity
	NegOne                    // Push Negative One
	Neq                       // Numeric Inequality
	Nil                       // Push Nil
	NoOp                      // Non-Operator
	Not                       // Boolean Negation
	One                       // Push One
	Panic                     // Abnormally Halt
	Pop                       // Discard Value
	PosInfinity               // Positive Infinity
	Resolve                   // Resolve Global Symbol
	RestArg                   // Retrieve Arguments Tail
	Return                    // Return Value
	ReturnFalse               // Return False
	ReturnNil                 // Return Nil
	ReturnTrue                // Return True
	Self                      // Push Current Function
	Store                     // Store Local
	Sub                       // Subtraction
	True                      // Push True
	Two                       // Push Two
	Zero                      // Push Zero
)

// Word makes Opcode a Coder
func (i Opcode) Word() Word {
	return Word(i)
}
