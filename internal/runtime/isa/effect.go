package isa

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/internal/debug"
)

type (
	// Effect captures how an Instruction impacts the state of the machine
	Effect struct {
		Operand ActOn // Describe element the operand acts on

		Pop  int // A fixed number of items to be popped from the stack
		Push int // A fixed number of items to be pushed onto the stack

		DPop   bool // Dynamic number of items to be popped (in operand)
		Ignore bool // Skip this instruction (ex: Labels and NoOps)
		Exit   bool // Results in a termination of the abstract machine
	}

	ActOn int
)

const (
	Nothing ActOn = iota
	Arguments
	Captured
	Constants
	Integer
	Labels
	Locals
	Stack
)

// ErrEffectNotDeclared is raised when an attempt to forcefully retrieve an
// Effect fails
var ErrEffectNotDeclared = errors.New("effect not declared for opcode")

// Effects is a lookup table of instruction effects
var Effects = map[Opcode]*Effect{
	// Ignored Opcodes
	Label: {Ignore: true, Operand: Labels},
	NoOp:  {Ignore: true},

	// Argument, Environment and Closure Operations
	Arg:        {Push: 1, Operand: Arguments},
	ArgsLen:    {Push: 1},
	ArgsPop:    {},
	ArgsPush:   {DPop: true, Operand: Arguments},
	ArgsRest:   {Push: 1, Operand: Arguments},
	Closure:    {Push: 1, Operand: Captured},
	EnvBind:    {Pop: 2},
	EnvPrivate: {Pop: 1},
	EnvPublic:  {Pop: 1},
	EnvValue:   {Pop: 1, Push: 1},

	// Reference and Register Operations
	Load:     {Push: 1, Operand: Locals},
	NewRef:   {Push: 1},
	RefBind:  {Pop: 2},
	RefValue: {Pop: 1, Push: 1},
	Store:    {Pop: 1, Operand: Locals},

	// Stack and Constant Operations
	Const: {Push: 1, Operand: Constants},
	Dup:   {Pop: 1, Push: 2},
	False: {Push: 1},
	Null:  {Push: 1},
	Pop:   {Pop: 1},
	Swap:  {Pop: 2, Push: 2},
	True:  {Push: 1},
	Zero:  {Push: 1},

	// Call Operations
	Call:     {Pop: 1, Push: 1, DPop: true, Operand: Stack},
	Call0:    {Pop: 1, Push: 1},
	Call1:    {Pop: 2, Push: 1},
	Call2:    {Pop: 3, Push: 1},
	Call3:    {Pop: 4, Push: 1},
	CallSelf: {Push: 1, DPop: true, Operand: Stack},
	CallWith: {Pop: 2, Push: 1},
	TailCall: {Pop: 1, DPop: true, Operand: Stack},
	TailClos: {Pop: 1, DPop: true, Operand: Stack},
	TailSelf: {DPop: true, Operand: Stack},

	// Control Flow Operations
	CondJump: {Pop: 1, Operand: Labels},
	Delay:    {Pop: 1, Push: 1},
	Jump:     {Operand: Labels},
	Panic:    {Pop: 1, Exit: true},
	RetFalse: {Exit: true},
	RetNull:  {Exit: true},
	RetTrue:  {Exit: true},
	Return:   {Pop: 1, Exit: true},

	// Sequence Operations
	Append:  {Pop: 2, Push: 1},
	Assoc:   {Pop: 2, Push: 1},
	Car:     {Pop: 1, Push: 1},
	Cdr:     {Pop: 1, Push: 1},
	Cons:    {Pop: 2, Push: 1},
	Dissoc:  {Pop: 2, Push: 1},
	Empty:   {Pop: 1, Push: 1},
	Get:     {Pop: 2, Push: 2},
	LazySeq: {Pop: 1, Push: 1},
	Length:  {Pop: 1, Push: 1},
	Nth:     {Pop: 2, Push: 2},
	Reverse: {Pop: 1, Push: 1},
	Vector:  {Push: 1, DPop: true, Operand: Stack},

	// Boolean Operations
	Eq:  {Pop: 2, Push: 1},
	Not: {Pop: 1, Push: 1},

	// Numeric Operations
	Add:    {Pop: 2, Push: 1},
	Div:    {Pop: 2, Push: 1},
	Mod:    {Pop: 2, Push: 1},
	Mul:    {Pop: 2, Push: 1},
	Neg:    {Pop: 1, Push: 1},
	NegInt: {Push: 1, Operand: Integer},
	NumEq:  {Pop: 2, Push: 1},
	NumGt:  {Pop: 2, Push: 1},
	NumGte: {Pop: 2, Push: 1},
	NumLt:  {Pop: 2, Push: 1},
	NumLte: {Pop: 2, Push: 1},
	PosInt: {Push: 1, Operand: Integer},
	Sub:    {Pop: 2, Push: 1},
}

func GetEffect(oc Opcode) (*Effect, error) {
	if effect, ok := Effects[oc]; ok {
		return effect, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrEffectNotDeclared, oc.String())
}

// MustGetEffect gives you effect information or explodes violently
func MustGetEffect(oc Opcode) *Effect {
	e, err := GetEffect(oc)
	if err != nil {
		panic(debug.ProgrammerErrorf("%w", err))
	}
	return e
}
