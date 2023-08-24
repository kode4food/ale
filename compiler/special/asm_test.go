package special_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestAsm(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define* test 
		  (lambda () (asm* one two add)))
		(test)
	`, I(3))
}
