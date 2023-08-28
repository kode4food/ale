package isa

import "fmt"

type (
	// Effect captures how an instruction impacts the state of the machine
	Effect struct {
		Operands []ActOn // Describe elements the operands act on

		Pop  int // Fixed number of items to be popped from the stack
		Push int // Fixed number of items to be pushed onto the stack
		DPop int // Dynamic number of items to be popped (first operand)

		Ignore bool // Skip this instruction (ex: Labels and NoOps)
		Exit   bool // Results in a termination of the VM
	}

	ActOn int
)

const (
	Nothing ActOn = iota
	Locals
	Constants
	Labels
	Arguments
	Stack
	Values
)

// Error messages
const (
	ErrEffectNotDeclared = "effect not declared for opcode: %s"
)

// Effects is a lookup table of instruction effects
var Effects = map[Opcode]*Effect{
	Add:        {Pop: 2, Push: 1},
	Arg:        {Push: 1, Operands: []ActOn{Arguments}},
	ArgLen:     {Push: 1},
	Bind:       {Pop: 2},
	BindRef:    {Pop: 2},
	Call:       {Pop: 1, Push: 1, DPop: 1, Operands: []ActOn{Stack}},
	Call0:      {Pop: 1, Push: 1},
	Call1:      {Pop: 2, Push: 1},
	Closure:    {Push: 1, Operands: []ActOn{Values}},
	CondJump:   {Pop: 1, Operands: []ActOn{Labels}},
	Const:      {Push: 1, Operands: []ActOn{Constants}},
	Deref:      {Pop: 1, Push: 1},
	Div:        {Pop: 2, Push: 1},
	Dup:        {Pop: 1, Push: 2},
	Eq:         {Pop: 2, Push: 1},
	False:      {Push: 1},
	Gt:         {Pop: 2, Push: 1},
	Gte:        {Pop: 2, Push: 1},
	Jump:       {Operands: []ActOn{Labels}},
	Label:      {Ignore: true, Operands: []ActOn{Labels}},
	Load:       {Push: 1, Operands: []ActOn{Locals}},
	Lt:         {Pop: 2, Push: 1},
	Lte:        {Pop: 2, Push: 1},
	MakeTruthy: {Pop: 1, Push: 1},
	Mod:        {Pop: 2, Push: 1},
	Mul:        {Pop: 2, Push: 1},
	Neg:        {Pop: 1, Push: 1},
	NegInf:     {Push: 1},
	NegOne:     {Push: 1},
	Neq:        {Pop: 2, Push: 1},
	NewRef:     {Push: 1},
	NoOp:       {Ignore: true},
	Not:        {Pop: 1, Push: 1},
	Nil:        {Push: 1},
	One:        {Push: 1},
	Panic:      {Pop: 1, Exit: true},
	Pop:        {Pop: 1},
	PosInf:     {Push: 1},
	Declare:    {Pop: 1},
	Private:    {Pop: 1},
	Resolve:    {Pop: 1, Push: 1},
	RestArg:    {Push: 1, Operands: []ActOn{Arguments}},
	RetFalse:   {Exit: true},
	RetNil:     {Exit: true},
	RetTrue:    {Exit: true},
	Return:     {Pop: 1, Exit: true},
	Self:       {Push: 1},
	Store:      {Pop: 1, Operands: []ActOn{Locals}},
	Sub:        {Pop: 2, Push: 1},
	TailCall:   {Pop: 1, DPop: 1, Operands: []ActOn{Stack}},
	True:       {Push: 1},
	Two:        {Push: 1},
	Zero:       {Push: 1},
}

// MustGetEffect gives you effect information or explodes violently
func MustGetEffect(oc Opcode) *Effect {
	if effect, ok := Effects[oc]; ok {
		return effect
	}
	panic(fmt.Errorf(ErrEffectNotDeclared, oc.String()))
}
