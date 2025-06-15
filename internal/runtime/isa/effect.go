package isa

import (
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
	Constants
	Integer
	Labels
	Locals
	Stack
	Values
)

// ErrEffectNotDeclared is raised when an attempt to forcefully retrieve an
// Effect fails
const ErrEffectNotDeclared = "effect not declared for opcode: %s"

// Effects is a lookup table of instruction effects
var Effects = map[Opcode]*Effect{
	Add:        {Pop: 2, Push: 1},
	Arg:        {Push: 1, Operand: Arguments},
	ArgsLen:    {Push: 1},
	ArgsPop:    {},
	ArgsPush:   {DPop: true, Operand: Arguments},
	ArgsRest:   {Push: 1, Operand: Arguments},
	Call:       {Pop: 1, Push: 1, DPop: true, Operand: Stack},
	Call0:      {Pop: 1, Push: 1},
	Call1:      {Pop: 2, Push: 1},
	Call2:      {Pop: 3, Push: 1},
	Call3:      {Pop: 4, Push: 1},
	CallSelf:   {Push: 1, DPop: true, Operand: Stack},
	CallWith:   {Pop: 2, Push: 1},
	Car:        {Pop: 1, Push: 1},
	Cdr:        {Pop: 1, Push: 1},
	Closure:    {Push: 1, Operand: Values},
	CondJump:   {Pop: 1, Operand: Labels},
	Cons:       {Pop: 2, Push: 1},
	Const:      {Push: 1, Operand: Constants},
	Div:        {Pop: 2, Push: 1},
	Dup:        {Pop: 1, Push: 2},
	Empty:      {Pop: 1, Push: 1},
	EnvBind:    {Pop: 2},
	EnvPrivate: {Pop: 1},
	EnvPublic:  {Pop: 1},
	EnvValue:   {Pop: 1, Push: 1},
	Eq:         {Pop: 2, Push: 1},
	False:      {Push: 1},
	Jump:       {Operand: Labels},
	Label:      {Ignore: true, Operand: Labels},
	Load:       {Push: 1, Operand: Locals},
	Mod:        {Pop: 2, Push: 1},
	Mul:        {Pop: 2, Push: 1},
	Neg:        {Pop: 1, Push: 1},
	NegInt:     {Push: 1, Operand: Integer},
	NewRef:     {Push: 1},
	NoOp:       {Ignore: true},
	Not:        {Pop: 1, Push: 1},
	Null:       {Push: 1},
	NumEq:      {Pop: 2, Push: 1},
	NumGt:      {Pop: 2, Push: 1},
	NumGte:     {Pop: 2, Push: 1},
	NumLt:      {Pop: 2, Push: 1},
	NumLte:     {Pop: 2, Push: 1},
	Panic:      {Pop: 1, Exit: true},
	Pop:        {Pop: 1},
	PosInt:     {Push: 1, Operand: Integer},
	RefBind:    {Pop: 2},
	RefValue:   {Pop: 1, Push: 1},
	RetFalse:   {Exit: true},
	RetNull:    {Exit: true},
	RetTrue:    {Exit: true},
	Return:     {Pop: 1, Exit: true},
	Store:      {Pop: 1, Operand: Locals},
	Sub:        {Pop: 2, Push: 1},
	TailCall:   {Pop: 1, DPop: true, Operand: Stack},
	TailClos:   {Pop: 1, DPop: true, Operand: Stack},
	TailSelf:   {DPop: true, Operand: Stack},
	True:       {Push: 1},
	Vector:     {Push: 1, DPop: true, Operand: Stack},
	Zero:       {Push: 1},
}

func GetEffect(oc Opcode) (*Effect, error) {
	if effect, ok := Effects[oc]; ok {
		return effect, nil
	}
	return nil, fmt.Errorf(ErrEffectNotDeclared, oc.String())
}

// MustGetEffect gives you effect information or explodes violently
func MustGetEffect(oc Opcode) *Effect {
	e, err := GetEffect(oc)
	if err != nil {
		panic(debug.ProgrammerErrorf("%w", err))
	}
	return e
}
