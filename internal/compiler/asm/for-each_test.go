package asm_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/asm"
)

func TestForEachParsingErrors(t *testing.T) {
	as := assert.New(t)

	as.ErrorWith(
		`(asm for-each)`,
		fmt.Errorf(asm.ErrExpectedType, "binding pair", ""),
	)

	as.ErrorWith(`
		(asm
			for-each (val)
				.eval val
			end)
	`, fmt.Errorf(asm.ErrExpectedType, "binding pair", "(val)"))

	as.ErrorWith(`
		(asm
			for-each [val]
				eval val
			end)
	`, fmt.Errorf(asm.ErrPairExpected, 1))

	as.ErrorWith(`
		(asm
			for-each [val body]
				eval val
			end)
	`, fmt.Errorf(asm.ErrUnexpectedParameter, "body"))

	as.ErrorWith(`
		(define t (special (body)
			for-each [val body]
				definitely-not-a-directive
			end))
	`, fmt.Errorf(asm.ErrUnknownDirective, "definitely-not-a-directive"))

	as.ErrorWith(`
		(define t (special (body)
				for-each [val body]
					eval val
				end))
		(t 1)
	`, fmt.Errorf(asm.ErrExpectedType, "sequence", "1"))
}
