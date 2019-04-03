package isa

import (
	"bytes"
	"fmt"
)

type analyzer struct {
	code    []Code
	maxSize int
	endSize int
	base    int
	pc      int
}

// Error messages
const (
	EffectNotDeclared = "effect not declared for opcode: %s"
)

// Verify checks an ISA code stream for validity. Specifically it will
// check that jumps do not target offsets outside of the instructions
// provided and that the stack is left in a consistent state upon exit
func Verify(code []Code) {
	verifyJumps(code)
	verifyStackSize(code)
}

// CalculateStackSize returns the maximum and final depths for the stack
// based on the instructions provided. If the final depth is non-zero,
// this is usually an indication that bad instructions were encoded
func CalculateStackSize(code []Code) (int, int) {
	a := &analyzer{
		code: code,
	}
	a.calculateStackSize()
	return a.maxSize, a.endSize
}

func verifyJumps(code []Code) {
	for i := 0; i < len(code); {
		oc := Opcode(code[i])
		effect := mustGetEffect(oc)
		if oc == CondJump || oc == Jump {
			off := int(code[i+1])
			if off < i || off > len(code) {
				fmt.Println(CodeToString(code, -1))
				panic(fmt.Sprintf("invalid jump offset: %d", off))
			}
		}
		i += effect.Size
	}
}

func verifyStackSize(code []Code) {
	a := &analyzer{
		code: code,
	}
	a.calculateStackSize()
	if a.endSize != 0 {
		fmt.Println(CodeToString(code, -1))
		panic(fmt.Sprintf("invalid stack end-state: %d", a.endSize))
	}
}

func (a *analyzer) calculateStackSize() {
	code := a.code
	codeLen := len(code)
	for a.pc < codeLen {
		oc := Opcode(code[a.pc])
		effect := mustGetEffect(oc)
		dPop := a.getStackChange(effect.DPop)
		dPush := a.getStackChange(effect.DPush)
		a.endSize += (effect.Push - effect.Pop) + (dPush - dPop)
		a.maxSize = maxInt(a.endSize, a.maxSize)
		if oc != CondJump || !a.calculateCondJump() {
			a.pc += effect.Size
		}
	}
}

func (a *analyzer) calculateIndex(code Code) int {
	return int(code) - a.base
}

func (a *analyzer) calculateCondJump() bool {
	effect := mustGetEffect(CondJump)
	thenIdx := a.calculateIndex(a.code[a.pc+1])
	if thenIdx <= a.pc {
		return false
	}

	before := coalesceInstructions(a.code[0:thenIdx])
	last := before[len(before)-1]
	if Opcode(last.code[0]) != Jump {
		return false
	}

	elseIdx := a.pc + effect.Size
	elseRes := a.calculateBranch(elseIdx, thenIdx)

	endIdx := a.calculateIndex(last.code[1])
	thenRes := a.calculateBranch(thenIdx, endIdx)

	if elseRes.endSize != thenRes.endSize {
		panic("branches should end in the same state")
	}

	a.endSize += elseRes.endSize
	a.pc = endIdx
	return true
}

func (a *analyzer) calculateBranch(start, end int) *analyzer {
	res := &analyzer{
		code:    a.code[start:end],
		maxSize: a.maxSize,
		endSize: a.endSize,
		base:    a.base + start,
	}
	res.calculateStackSize()
	a.maxSize = maxInt(a.maxSize, res.maxSize)
	return res
}

func (a *analyzer) getStackChange(count int) int {
	if count > 0 {
		return int(a.code[a.pc+count])
	}
	return 0
}

type inst struct {
	pc   int
	code []Code
}

func coalesceInstructions(code []Code) []inst {
	var res []inst
	for i := 0; i < len(code); {
		oc := Opcode(code[i])
		effect := mustGetEffect(oc)
		size := effect.Size
		elem := make([]Code, size)
		copy(elem, code[i:])
		res = append(res, inst{
			pc:   i,
			code: elem,
		})
		i += size
	}
	return res
}

func mustGetEffect(oc Opcode) *Effect {
	if effect, ok := Effects[oc]; ok {
		return effect
	}
	panic(fmt.Sprintf(EffectNotDeclared, oc.String()))
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

// CodeToString converts a code array to a string form for debuggign
func CodeToString(code []Code, PC int) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("PC = %d\n", PC))
	inst := coalesceInstructions(code)
	for _, line := range inst {
		buf.WriteString(fmt.Sprintf("%d: ", line.pc))
		oc := Opcode(line.code[0])
		buf.WriteString(fmt.Sprintf("%s", oc.String()))
		if len(line.code) > 1 {
			buf.WriteString("(")
			for _, operand := range line.code[1:] {
				if len(line.code) > 2 {
					buf.WriteString(", ")
				}
				buf.WriteString(fmt.Sprintf("%d", operand))
			}
			buf.WriteString(")")
		}
		buf.WriteString("\n")
	}
	return buf.String()
}
