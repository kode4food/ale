package asm_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/core/asm"
)

func TestUnknownDirective(t *testing.T) {
	as := assert.New(t)
	as.ErrorWith(`
		(asm*
			.push-locals
			definitely-not-a-directive
			.pop-locals)
	`, fmt.Errorf(asm.ErrUnknownDirective, "definitely-not-a-directive"))
}

func TestBadArgs(t *testing.T) {
	as := assert.New(t)
	as.ErrorWith(`
		(asm* .const)
	`, fmt.Errorf(asm.ErrTooFewArguments, ".const", 1))

	as.ErrorWith(`
		(asm* .resolve 99)
	`, fmt.Errorf(asm.ErrExpectedType, "symbol", "99"))
}

func TestUnknownForm(t *testing.T) {
	as := assert.New(t)
	as.ErrorWith(`
		(asm*
			.push-locals
			{:not "valid"}
			.pop-locals)
	`, fmt.Errorf(asm.ErrUnexpectedForm, `{:not "valid"}`))
}

func TestBadBlocks(t *testing.T) {
	as := assert.New(t)
	as.ErrorWith(`
	  (special* body
		  .for-each [val body]
			  .eval val))
	`, fmt.Errorf(asm.ErrExpectedEndOfBlock))

	as.ErrorWith(`
	  (special* body
		  .for-each [val body]
			  definitely-not-a-directive
		  .end))
	`, fmt.Errorf(asm.ErrUnknownDirective, "definitely-not-a-directive"))

	as.ErrorWith(`
	  (special* body
		  .for-each [val body]
			  .eval val
              null
		  .end
		  definitely-not-a-directive))
	`, fmt.Errorf(asm.ErrUnknownDirective, "definitely-not-a-directive"))
}
